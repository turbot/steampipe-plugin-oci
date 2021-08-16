package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/functions"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableFunctionsApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_functions_application",
		Description: "OCI Functions Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getFunctionsApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listFunctionsApplications,
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
				Description: "The OCID of the application.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "syslog_url",
				Description: "A syslog URL to which to send all function logs. Supports tcp, udp, and tcp+tls.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the application was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the application was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			//json fields
			{
				Name:        "config",
				Description: "Application configuration for functions in this application.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFunctionsApplication,
			},
			{
				Name:        "subnet_ids",
				Description: "The OCIDs of the subnets in which to run functions in the application.",
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
				Transform:   transform.From(applicationTags),
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

func listFunctionsApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listFunctionsApplications", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := functionsManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := functions.ListApplicationsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.FunctionsManagementClient.ListApplications(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, application := range response.Items {
			d.StreamListItem(ctx, application)
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

func getFunctionsApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFunctionsApplication")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFunctionsApplication", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		application := h.Item.(functions.ApplicationSummary)
		id = *application.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty application id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := functionsManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := functions.GetApplicationRequest{
		ApplicationId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.FunctionsManagementClient.GetApplication(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Application, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func applicationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	freeformTags := applicationFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := applicationDefinedTags(d.HydrateItem)

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

func applicationFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case functions.Application:
		return item.FreeformTags
	case functions.ApplicationSummary:
		return item.FreeformTags
	}
	return nil
}

func applicationDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case functions.Application:
		return item.DefinedTags
	case functions.ApplicationSummary:
		return item.DefinedTags
	}
	return nil
}
