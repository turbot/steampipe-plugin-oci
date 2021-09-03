package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/apigateway"
	oci_common "github.com/oracle/oci-go-sdk/v44/common"
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
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getApiGatewayApi,
		},
		List: &plugin.ListConfig{
			Hydrate: listApiGatewayApis,
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
				Name:        "lifecycle_state",
				Description: "The current state of the API.",
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
				Name:        "lifecycle_details",
				Description: "A message describing the current lifecycleState.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "specification_type",
				Description: "Type of API Specification file.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The time this resource was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			// json fields
			{
				Name:        "validation_results",
				Description: "Status of each feature available from the API.",
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
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listApiGatewayApis(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

func getApiGatewayApi(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getApiGatewayApi")
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty api id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
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
	freeformTags := apiGatewayApiFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := apiGatewayApiDefinedTags(d.HydrateItem)

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

func apiGatewayApiFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case apigateway.Api:
		return item.FreeformTags
	case apigateway.ApiSummary:
		return item.FreeformTags
	}
	return nil
}

func apiGatewayApiDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case apigateway.Api:
		return item.DefinedTags
	case apigateway.ApiSummary:
		return item.DefinedTags
	}
	return nil
}
