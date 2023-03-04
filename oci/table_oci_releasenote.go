package oci

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableReleaseNote(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_releasenote",
		Description: "OCI Release Note",
		List: &plugin.ListConfig{
			Hydrate: listReleaseNotes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "service",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "title",
				Description: "The title of the release note.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Title"),
			},
			{
				Name:        "summary",
				Description: "Short description of the release note.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Summary"),
			},
			{
				Name:        "release_date",
				Description: "Publication date for release note.",
				Type:        proto.ColumnType_TIMESTAMP,
				Default:     false,
				Transform:   transform.FromField("ReleaseDate"),
			},
			{
				Name:        "service",
				Description: "Primary OCI Service involved in this release.",
				Type:        proto.ColumnType_STRING,
				Default:     false,
				Transform:   transform.FromField("Service"),
			},
			{
				Name:        "all_services",
				Description: "Array of all OCI Services related to this release.",
				Type:        proto.ColumnType_JSON,
				Default:     false,
				Transform:   transform.From(transformAllServices),
			},
			{
				Name:        "url",
				Description: "URL to webpage with details on the release note.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Url"),
			},

			// Standard Steampipe columns - N/A for this table

			// Standard OCI columns - N/A for this table
		},
	}
}

type ReleaseNote struct {
	Title       string
	Service     string
	AllServices []string
	URL         string
	ReleaseDate *time.Time
	Summary     string
}

//// LIST FUNCTION

// TODO release notes for a specific service are also available on service specific pages, such as https://docs.oracle.com/en-us/iaas/releasenotes/services/api-gateway/
// TODO the url is composed of https://docs.oracle.com/en-us/iaas/releasenotes/services/<url encoded name of the sercice>/
// TODO when the where clause contains a reference to the service, the data can be fetched from this far smaller subset of focused data

func listReleaseNotes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("listReleaseNotes")

	equalQuals := d.KeyColumnQuals
	serviceFilter := strings.ToLower(strings.ReplaceAll(equalQuals["service"].GetStringValue(), " ", "-"))

	// TODO: some services have a special mapping to the corresponding URL; for example:
	// application-performance-monitoring => apm
	// artifact-registry => artifacts
	// block-volume => blockvolume
	if serviceFilter == "application-performance-monitoring" {
		serviceFilter = "apm"
	}
	if serviceFilter == "artifact-registry" {
		serviceFilter = "artifacts"
	}
	if serviceFilter == "block-volume" {
		serviceFilter = "blockvolume"
	}
	if serviceFilter == "file-storage" {
		serviceFilter = "filestorage"
	}

	logger.Debug("listReleaseNotes, filter on service " + serviceFilter)

	pagesLeft := true
	currentPage := 0
	for pagesLeft {
		currentPage++
		releaseNotes, numberOfPages, err := getOCIReleaseNotes(serviceFilter, currentPage)
		if err != nil {
			logger.Error("oci_releasenote.listReleaseNotes", "scrape data error", err)
			return nil, err
		}
		for _, releaseNote := range releaseNotes {
			d.StreamListItem(ctx, releaseNote)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		// no futher data retrieval when the current page is the last page or when the service filter was defined
		pagesLeft = currentPage < numberOfPages && serviceFilter == ""
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

// RELEASE_DATE_FORMAT one of two date formats date used in release notes
const RELEASE_DATE_FORMAT = "Jan. 2, 2006"

// ALTERNATIVE_DATE_FORMAT the other one of the two date formats date used in release notes
const ALTERNATIVE_DATE_FORMAT = "January 2, 2006"

const ALLOWED_DOMAIN = "docs.oracle.com"

// note: the most recent 50 release notes are published as RSS feed at: https://docs.oracle.com/en-us/iaas/releasenotes/feed/
// using this feed is (far) less brittle than the screenscraping approach

// the url where the web pages with release notes are available
const RELEASE_NOTES_BASE_URL = "https://docs.oracle.com/en-us/iaas/releasenotes/"
const BASE_DOCUMENTATION_URL = "https://docs.oracle.com"

func getOCIReleaseNotes(service string, page int) (releaseNotes []ReleaseNote, numberOfPages int, err error) {
	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAIN),
	)
	releaseNotes = make([]ReleaseNote, 0, 50)
	c.OnHTML(".uk-article", func(e *colly.HTMLElement) {
		article := ReleaseNote{}
		article.Title = e.ChildText(" h3 a[href]")
		article.URL = BASE_DOCUMENTATION_URL + e.ChildAttr("h3 a", "href")
		services := make([]string, 0, 20)
		e.ForEach("ul li:nth-child(1) a[href]", func(index int, a *colly.HTMLElement) {
			services = append(services, a.Text)
		})

		article.AllServices = make([]string, 0, len(services))
		copy(article.AllServices, services[:])
		article.Service = services[0]
		article.AllServices = services
		article.Summary = e.ChildText(" div.uk-panel p ")
		dateString := e.ChildText(" ul  li:nth-child(2)")
		_, after, _ := strings.Cut(dateString, "Release Date: ")
		releaseDateString := after
		releaseDate, error := time.Parse(RELEASE_DATE_FORMAT, releaseDateString)
		if error != nil {
			// try to parse the date string with the alternative formay
			altReleaseDate, _ := time.Parse(ALTERNATIVE_DATE_FORMAT, releaseDateString)
			releaseDate = altReleaseDate
		}
		article.ReleaseDate = &releaseDate
		releaseNotes = append(releaseNotes, article)
	})
	c.OnHTML(".page-current", func(e *colly.HTMLElement) {
		numberOfPages, err = getNumberOfPages(e.Text)
		if err != nil {
			return
		}
	})

	if service == "" {
		err = c.Visit(RELEASE_NOTES_BASE_URL + "?page=" + strconv.Itoa(page))
		if (err != nil) {
			return
		}

	}
	if service != "" {
		// derive service specific page with release notes!
		// https://docs.oracle.com/en-us/iaas/releasenotes/services/speech/

		err = c.Visit(RELEASE_NOTES_BASE_URL + "/services/" + service + "/")
		if (err != nil) {
			return
		}
	}
	return
}

func getNumberOfPages(text string) (int, error) {
	// number of pages is scraped from page <span class="page-current">Page 1 of 36</span>
	// the text content of the span is passed to this function

	_, after, _ := strings.Cut(text, " of ")
	before, _, _ := strings.Cut(after, "\n")

	numberOfPages, err := strconv.Atoi(before)
	if err != nil {
		//	fmt.Println("Error converting string to integer:"+before2, err)
		return 0, err
	}
	return numberOfPages, nil
}

// produce a valid JSON representation of the slice of strings in field AllServices
func transformAllServices(_ context.Context, d *transform.TransformData) (interface{}, error) {
	releaseNote := d.HydrateItem.(ReleaseNote)
	allServices := releaseNote.AllServices
	allServicesJSON, err := json.Marshal(allServices)
	if err != nil {
		return nil, err
	}

	return string(allServicesJSON), nil
}
