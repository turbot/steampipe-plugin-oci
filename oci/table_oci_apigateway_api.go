package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/apigateway"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableApiGatewayApi(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_apigateway_api",
		Description: "OCI Apigateway Api",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getApigatewayApi,
		},
		List: &plugin.ListConfig{
			Hydrate: listApigatewayApis,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time this resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "time_updated",
				Description: "The OCID of dedicated VM host.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current lifecycleState",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "specification_type",
				Description: "Type of API Specification file.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "validation_results",
				Description: "Status of each feature available from the API.",
				Type:        proto.ColumnType_JSON,
			},

			// json fields

			// tags
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a namespace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(apiTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(ociRegionName),
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

//// LIST FUNCTION

func listApigatewayApis(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	// Create Session
	session, err := apiGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := apigateway.ListApisRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ApiGatewayClient.ListApis(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, api := range response.Items {
			d.StreamListItem(ctx, api)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getApigatewayApi(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApigatewayApi")
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	var id string
	if h.Item != nil {
		id = apiId(h.Item)
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// Create Session
	session, err := apiGatewayService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := apigateway.GetApiRequest{
		ApiId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ApiGatewayClient.GetApi(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Api, nil
}

// //// TRANSFORM FUNCTION

// // Priority order for tags
// // 1. System Tags
// // 2. Defined Tags
// // 3. Free-form tags
func apiTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := apigatewayApiFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := apigatewayApiDefinedTags(d.HydrateItem)

	if definedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range definedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

func apiId(item interface{}) string {
	switch item.(type) {
	case apigateway.Api:
		return *item.(apigateway.Api).Id
	case apigateway.ApiSummary:
		return *item.(apigateway.ApiSummary).Id
	}
	return ""
}

func apigatewayApiFreeformTags(item interface{}) map[string]string {
	switch item.(type) {
	case apigateway.Api:
		return item.(apigateway.Api).FreeformTags
	case apigateway.ApiSummary:
		return item.(apigateway.ApiSummary).FreeformTags
	}
	return nil
}

func apigatewayApiDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item.(type) {
	case apigateway.Api:
		return item.(apigateway.Api).DefinedTags
	case apigateway.ApiCollection:
		return item.(apigateway.ApiSummary).DefinedTags
	}
	return nil
}
