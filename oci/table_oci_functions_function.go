package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/functions"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableFunctionsFunction(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_functions_function",
		Description: "OCI Functions Function",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getFunction,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listFunctionsApplications,
			Hydrate:       listFunctions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The display name of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the function.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "application_id",
				Description: "The OCID of the application the function belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ApplicationId"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the function.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image",
				Description: "The qualified name of the Docker image to use in the function, including the image tag. The image should be in the OCI Registry that is in the same region as the function itself.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_digest",
				Description: "The image digest for the version of the image that will be pulled when invoking this function. If no value is specified, the digest currently associated with the image in the OCI Registry will be used.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "invoke_endpoint",
				Description: "The base https invoke URL to set on a client in order to invoke a function. This URL will never change over the lifetime of the function and can be cached.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "memory_in_mbs",
				Description: "Maximum usable memory for the function (MiB).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MemoryInMBs"),
			},
			{
				Name:        "time_created",
				Description: "The time the function was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the function was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "timeout_in_seconds",
				Description: "Timeout for executions of the function. Value in seconds.",
				Type:        proto.ColumnType_INT,
			},

			// json fields
			{
				Name:        "config",
				Description: "The function configuration. Overrides application configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunction,
			},
			{
				Name:        "trace_config",
				Description: "The trace configuration of the function.",
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
				Transform:   transform.From(functionTags),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listFunctions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Trace("listFunctions", "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	var applicationId string

	if equalQuals["application_id"] != nil {
		applicationId = *types.String(equalQuals["application_id"].GetStringValue())
	} else {
		applicationId = *h.Item.(functions.ApplicationSummary).Id
	}

	// handle empty application id in list call
	if applicationId == "" {
		return nil, nil
	}

	// Create Session
	session, err := functionsManagementService(ctx, d, region)
	if err != nil {
		logger.Error("listFunctions", "error_functionsManagementService", err)
		return nil, err
	}
	
	// Build request parameters
	request := buildFunctionsFilters(equalQuals)
	request.ApplicationId = &applicationId
	request.Limit = types.Int(50)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.FunctionsManagementClient.ListFunctions(ctx, request)
		if err != nil {
			logger.Error("listFunctions", "error_ListFunctions", err)
			return nil, err
		}

		for _, item := range response.Items {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
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

func getFunction(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunction")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFunction", "Compartment", compartment, "OCI_REGION", region)

	var functionId string
	if h.Item != nil {
		functionId = *h.Item.(functions.FunctionSummary).Id
	} else {
		functionId = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty function id in get call
	if functionId == "" {
		return nil, nil
	}

	// Create Session
	session, err := functionsManagementService(ctx, d, region)
	if err != nil {
		logger.Error("getFunction", "error_functionsManagementService", err)
		return nil, err
	}

	request := functions.GetFunctionRequest{
		FunctionId: types.String(functionId),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.FunctionsManagementClient.GetFunction(ctx, request)
	if err != nil {
		logger.Error("getFunction", "error_GetFunction", err)
		return nil, err
	}

	return response.Function, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func functionTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	freeformTags := functionFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := functionDefinedTags(d.HydrateItem)

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

func functionFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case functions.Function:
		return item.FreeformTags
	case functions.FunctionSummary:
		return item.FreeformTags
	}
	return nil
}

func functionDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case functions.Function:
		return item.DefinedTags
	case functions.FunctionSummary:
		return item.DefinedTags
	}
	return nil
}

// Build additional filters
func buildFunctionsFilters(equalQuals plugin.KeyColumnEqualsQualMap) functions.ListFunctionsRequest {
	request := functions.ListFunctionsRequest{}

	if equalQuals["id"] != nil {
		request.Id = types.String(equalQuals["id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = functions.FunctionLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
