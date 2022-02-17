package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/logging"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableLoggingLog(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_logging_log",
		Description: "OCI Logging Log",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "log_group_id"}),
			Hydrate:    getLoggingLog,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listLoggingLogGroups,
			Hydrate:       listLoggingLogs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "log_group_id",
					Require: plugin.Optional,
				},
				{
					Name:    "log_type",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "The OCID of the log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "log_group_id",
				Description: "The OCID of the log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The log object state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "log_type",
				Description: "The logType that the log object is for, whether custom or service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_enabled",
				Description: "Whether or not this resource is currently enabled.",
				Type:        proto.ColumnType_BOOL,
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
			{
				Name:        "retention_duration",
				Description: "Log retention duration in 30-day increments (30, 60, 90 and so on).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "configuration",
				Description: "Log object configuration.",
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
				Transform:   transform.From(logTags),
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

func listLoggingLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listLoggingLogs", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Create Session
	session, err := loggingManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	logGroupId := *h.Item.(logging.LogGroupSummary).Id
	
	
	// Build request parameters
	request := buildLoggingLogFilters(equalQuals)
	request.LogGroupId = types.String(logGroupId)
	request.Limit = types.Int(1000)
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
		response, err := session.LoggingManagementClient.ListLogs(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, log := range response.Items {
			d.StreamListItem(ctx, log)

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

func getLoggingLog(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLoggingLog")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getLoggingLog", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()
	logGroupId := d.KeyColumnQuals["log_group_id"].GetStringValue()

	// handle empty log and log group id in get call
	if id == "" || logGroupId == "" {
		return nil, nil
	}

	// Create Session
	session, err := loggingManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := logging.GetLogRequest{
		LogGroupId: types.String(logGroupId),
		LogId:      types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.LoggingManagementClient.GetLog(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Log, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 2. Defined Tags
// 3. Free-form tags
func logTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := logFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := logDefinedTags(d.HydrateItem)

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

func logFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case logging.Log:
		return item.FreeformTags
	case logging.LogSummary:
		return item.FreeformTags
	}
	return nil
}

func logDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case logging.Log:
		return item.DefinedTags
	case logging.LogSummary:
		return item.DefinedTags
	}
	return nil
}

// Build additional filters
func buildLoggingLogFilters(equalQuals plugin.KeyColumnEqualsQualMap) logging.ListLogsRequest {
	request := logging.ListLogsRequest{}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = logging.ListLogsLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}
	if equalQuals["log_group_id"] != nil {
		request.LogGroupId = types.String(equalQuals["log_group_id"].GetStringValue())
	}
	if equalQuals["log_type"] != nil {
		request.LogType = logging.ListLogsLogTypeEnum(equalQuals["log_type"].GetStringValue())
	}
	if equalQuals["name"] != nil {
		request.DisplayName = types.String(equalQuals["name"].GetStringValue())
	}

	return request
}
