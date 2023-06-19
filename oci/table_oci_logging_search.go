package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/loggingsearch"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableLoggingSearch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_logging_search",
		Description: "OCI Logging Search",
		List: &plugin.ListConfig{
			Hydrate: listLoggingSearch,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "start_time",
					Require: plugin.Required,
				},
				{
					Name:    "end_time",
					Require: plugin.Required,
				},
				{
					Name:    "log_group_id",
					Require: plugin.Required,
				},
				{
					Name:    "log_id",
					Require: plugin.Required,
				},
				{
					Name:    "level",
					Require: plugin.Required,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "start_time",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("start_time"),
			},
			{
				Name:        "end_time",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromQual("end_time"),
			},
			{
				Name:        "log_group_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_group_id"),
			},
			{
				Name:        "log_id",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_id"),
			},
			{
				Name:        "level",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("level"),
			},
			{
				Name:        "summary",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "results",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "fields",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns

			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				//Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
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

func listLoggingSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listLoggingSearch", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := loggingSearchService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	request := loggingsearch.SearchLogsRequest{}
	request.SearchLogsDetails.IsReturnFieldInfo = types.Bool(true)
	request.SearchLogsDetails.TimeStart = &common.SDKTime{Time: d.KeyColumnQuals["start_time"].GetTimestampValue().AsTime()}
	request.SearchLogsDetails.TimeEnd = &common.SDKTime{Time: d.KeyColumnQuals["end_time"].GetTimestampValue().AsTime()}
	log_group_id := d.KeyColumnQualString("log_group_id")
	log_id := d.KeyColumnQualString("log_id")
	level := d.KeyColumnQualString("level")
	searchQuery := "search " + compartment + "/" + log_group_id + "/" + log_id + "| where level = " + level + ";"

	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}
	request.SearchLogsDetails.SearchQuery = types.String(searchQuery)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.LoggingSearchClient.SearchLogs(ctx, request)
		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, response.SearchResponse)

		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}
