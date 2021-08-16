package oci

import (
	"context"
	"errors"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/resourcesearch"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableResourceSearch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_resource_search",
		Description: "OCI Resource Search",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"query", "text"}),
			Hydrate:    listResourceSearch,
		},
		GetMatrixItem: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "identifier",
				Description: "The unique identifier for this particular resource, usually an OCID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "The display name (or name) of this resource, if one exists.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain where this resource exists, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_type",
				Description: "The resource type name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The lifecycle state of this resource, if applicable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time that this resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "query",
				Description: "The query based on which the search was done.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "text",
				Description: "The freeText based on which the search was done.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_context",
				Description: "Additional identifiers to use together in a Get request for a specified resource, only required for resource types that explicitly cannot be retrieved by using a single identifier, such as the resource's OCID.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "search_context",
				Description: "SearchContext Contains search context, such as highlighting, for found resources.",
				Type:        proto.ColumnType_JSON,
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "system_tags",
				Description: "System tags associated with this resource.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(resourceSearchTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "search_region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Region"),
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

type searchInfo struct {
	resourcesearch.ResourceSummary
	Query  string
	Region string
	Text   string
}

//// LIST FUNCTION

func listResourceSearch(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listResourceSearch", "OCI_REGION", region)

	query := d.KeyColumnQuals["query"].GetStringValue()
	text := d.KeyColumnQuals["text"].GetStringValue()

	// handle empty query and text in list call
	if query == "" && text == "" {
		return nil, nil
	}

	if query != "" && text != "" {
		return nil, errors.New("please provide either query or text")
	}

	// Create Session
	session, err := resourceSearchService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	if query != "" {
		request := resourcesearch.SearchResourcesRequest{
			SearchDetails: resourcesearch.StructuredSearchDetails{
				Query: common.String(query),
			},
			RequestMetadata: common.RequestMetadata{
				RetryPolicy: getDefaultRetryPolicy(),
			},
		}

		pagesLeft := true
		for pagesLeft {
			response, err := session.ResourceSearchClient.SearchResources(ctx, request)
			if err != nil {
				return nil, err
			}

			for _, resource := range response.Items {
				d.StreamListItem(ctx, searchInfo{resource, query, region, ""})
			}
			if response.OpcNextPage != nil {
				request.Page = response.OpcNextPage
			} else {
				pagesLeft = false
			}
		}
	}

	if text != "" {
		request := resourcesearch.SearchResourcesRequest{
			SearchDetails: resourcesearch.FreeTextSearchDetails{
				Text: common.String(text),
			},
			RequestMetadata: common.RequestMetadata{
				RetryPolicy: getDefaultRetryPolicy(),
			},
		}
		pagesLeft := true
		for pagesLeft {
			response, err := session.ResourceSearchClient.SearchResources(ctx, request)
			if err != nil {
				return nil, err
			}

			for _, resource := range response.Items {
				d.StreamListItem(ctx, searchInfo{resource, "", region, text})
			}
			if response.OpcNextPage != nil {
				request.Page = response.OpcNextPage
			} else {
				pagesLeft = false
			}
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func resourceSearchTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	resourceSearch := d.HydrateItem.(searchInfo)

	var tags map[string]interface{}

	if resourceSearch.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range resourceSearch.FreeformTags {
			tags[k] = v
		}
	}

	if resourceSearch.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range resourceSearch.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	if resourceSearch.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range resourceSearch.SystemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
