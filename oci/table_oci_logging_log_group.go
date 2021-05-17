package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/logging"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableLoggingLogGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_logging_log_group",
		Description: "OCI Logging Log Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getLoggingLogGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listLoggingLogGroups,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The log group object state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description for this log group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time the resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_last_modified",
				Description: "Time the resource was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeLastModified.Time"),
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
				Transform:   transform.From(logGroupTags),
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

func listLoggingLogGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listLoggingLogGroups", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := loggingManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := logging.ListLogGroupsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.LoggingManagementClient.ListLogGroups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, logGroup := range response.Items {
			d.StreamListItem(ctx, logGroup)
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

func getLoggingLogGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingLogGroup")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getLoggingLogGroup", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty internet gateway id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := loggingManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := logging.GetLogGroupRequest{
		LogGroupId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.LoggingManagementClient.GetLogGroup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.LogGroup, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func logGroupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := logGroupFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := logGroupDefinedTags(d.HydrateItem)

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

func logGroupFreeformTags(item interface{}) map[string]string {
	switch item.(type) {
	case logging.LogGroup:
		return item.(logging.LogGroup).FreeformTags
	case logging.LogGroupSummary:
		return item.(logging.LogGroupSummary).FreeformTags
	}
	return nil
}

func logGroupDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item.(type) {
	case logging.LogGroup:
		return item.(logging.LogGroup).DefinedTags
	case logging.LogGroupSummary:
		return item.(logging.LogGroupSummary).DefinedTags
	}
	return nil
}
