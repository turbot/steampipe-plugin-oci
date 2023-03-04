package oci

import (
	"context"
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

// TODO define and populate column to provide a list of all services references by a release note

func listReleaseNotes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("listReleaseNotes")

	//limit := d.QueryContext.Limit

	pagesLeft := true
	currentPage := 0
	for pagesLeft {
		currentPage++
		releaseNotes, numberOfPages, err := getOCIReleaseNotes(currentPage)
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
		pagesLeft = currentPage < numberOfPages
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

func getOCIReleaseNotes(page int) (releaseNotes []ReleaseNote, numberOfPages int, err error) {
	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAIN),
	)
	releaseNotes = make([]ReleaseNote, 0, 50)
	c.OnHTML(".uk-article", func(e *colly.HTMLElement) {
		article := ReleaseNote{}
		article.Title = e.ChildText(" h3 a[href]")
		article.URL = BASE_DOCUMENTATION_URL + e.ChildAttr("h3 a", "href")
		services := strings.Split(e.ChildText(" ul  li:nth-child(1) a "), ",")
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
	c.Visit(RELEASE_NOTES_BASE_URL + "?page=" + strconv.Itoa(page))
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
