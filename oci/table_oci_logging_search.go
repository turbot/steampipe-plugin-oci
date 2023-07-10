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
			Hydrate:           listLoggingSearch,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "start_time",
					Require: plugin.Optional,
				},
				{
					Name:    "end_time",
					Require: plugin.Optional,
				},
				{
					Name:    "log_group_name",
					Require: plugin.Optional,
				},
				{
					Name:    "log_name",
					Require: plugin.Optional,
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
				Name:        "log_group_name",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_group_name"),
			},
			{
				Name:        "log_name",
				Description: "",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_name"),
			},
			{
				Name:        "datetime",
				Description: "",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Data.datetime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "log_content",
				Description: "",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Data.logContent"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
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

type LoggingSearch struct {
	Data          interface{}
	CompartmentId string
	Region        string
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
	request.SearchLogsDetails.TimeStart = &common.SDKTime{Time: d.KeyColumnQuals["start_time"].GetTimestampValue().AsTime()}
	request.SearchLogsDetails.TimeEnd = &common.SDKTime{Time: d.KeyColumnQuals["end_time"].GetTimestampValue().AsTime()}
	log_group_name := d.KeyColumnQualString("log_group_name")
	log_name := d.KeyColumnQualString("log_name")

	// prepare the query
	query := compartment
	if log_group_name != "" {
		query = query + "/" + log_group_name
	}
	if log_name != "" {
		query = query + "/" + log_name
	}
	searchQuery := "search \"" + query + "\""
	request.SearchLogsDetails.SearchQuery = types.String(searchQuery)

	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

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
		for _, result := range response.SearchResponse.Results {
			d.StreamListItem(ctx, LoggingSearch{*result.Data, compartment, region})

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
