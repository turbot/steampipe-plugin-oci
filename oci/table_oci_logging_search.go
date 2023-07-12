package oci

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/loggingsearch"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
					Name:       "timestamp",
					Operators:  []string{">", ">=", "=", "<", "<="},
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "log_group_name",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "log_name",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:       "search_query",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "log_content_id",
				Description: "The log content id.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Data.logContent.id"),
			},
			{
				Name:        "log_content_source",
				Description: "The log content source.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Data.logContent.source"),
			},
			{
				Name:        "log_content_type",
				Description: "The log content type.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Data.logContent.type"),
			},
			{
				Name:        "timestamp",
				Description: "Represents the timestamp of a log entry.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Data.datetime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "log_group_name",
				Description: "Specifies the name of the log group to which the log entry belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_group_name"),
			},
			{
				Name:        "log_name",
				Description: "Indicates the name of the log within the log group. It helps to identify the specific log within a log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_name"),
			},
			{
				Name:        "search_query",
				Description: "Stores the search query associated with the log entry. It represents the criteria used to filter and retrieve specific logs from the log group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("search_query"),
			},
			{
				Name:        "log_content",
				Description: "Stores the actual content of the log entry in JSON format. It contains the detailed information about the log event, such as log message, metadata, and any additional structured data.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Data.logContent"),
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Data.logContent.id"),
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
				Transform:   transform.FromField("Data.logContent.oracle.compartmentid"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Data.logContent.oracle.tenantid"),
			},
		},
	}
}

type LoggingSearch struct {
	Data   interface{}
	Region string
}

//// LIST FUNCTION

func listLoggingSearch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listLoggingSearch", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := loggingSearchService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := loggingsearch.SearchLogsRequest{}

	//set the start and end time based on the provided timestamp
	var timeStart, timeEnd *common.SDKTime
	if d.Quals["timestamp"] != nil {
		for _, q := range d.Quals["timestamp"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case "=":
				timeStart = &common.SDKTime{Time: timestamp}
				timeEnd = &common.SDKTime{Time: timestamp}
			case ">=", ">":
				timeStart = &common.SDKTime{Time: timestamp}
			case "<", "<=":
				timeEnd = &common.SDKTime{Time: timestamp}
			}
		}
	}
	if timeStart == nil {
		timeStart = &common.SDKTime{Time: (time.Now().AddDate(0, 0, -1))}
	}
	if timeEnd == nil {
		timeEnd = &common.SDKTime{Time: (time.Now())}
	}
	request.SearchLogsDetails.TimeStart = timeStart
	request.SearchLogsDetails.TimeEnd = timeEnd
	var searchQuery string
	if d.EqualsQualString("search_query") == "" {
		log_group_name := d.EqualsQualString("log_group_name")
		log_name := d.EqualsQualString("log_name")

		// prepare the query
		query := compartment
		if log_group_name != "" {
			query = query + "/" + log_group_name
		}
		if log_name != "" {
			query = query + "/" + log_name
		}
		searchQuery = "search \"" + query + "\""
	} else {
		searchQuery = d.EqualsQualString("search_query")
	}
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
			d.StreamListItem(ctx, LoggingSearch{*result.Data, region})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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
