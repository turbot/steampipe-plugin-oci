package oci

import (
	"context"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v44/common"
	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/monitoring"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// append the common metric columns onto the column list
func MonitoringMetricColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonMonitoringMetricColumns()...)
}

func commonMonitoringMetricColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "metric_name",
			Description: "The name of the metric.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "namespace",
			Description: "The metric namespace.",
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "average",
			Description: "The average of the metric values that correspond to the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "maximum",
			Description: "The maximum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "minimum",
			Description: "The minimum metric value for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sample_count",
			Description: "The number of metric values that contributed to the aggregate value of this data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "sum",
			Description: "The sum of the metric values for the data point.",
			Type:        proto.ColumnType_DOUBLE,
		},
		{
			Name:        "unit",
			Description: "The standard unit for the data point.",
			Type:        proto.ColumnType_STRING,
			Transform:   transform.FromField("Metadata.unit"),
		},
		{
			Name:        "timestamp",
			Description: "The time stamp used for the data point.",
			Type:        proto.ColumnType_TIMESTAMP,
		},
		{
			Name:        "region",
			Description: ColumnDescriptionRegion,
			Type:        proto.ColumnType_STRING,
		},
		{
			Name:        "compartment_id",
			Description: "The ID of the compartment.",
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
	}
}

type MonitoringMetricRow struct {
	// The compartment id
	CompartmentId *string

	// The (single) metric Dimension name
	DimensionName *string

	// The value for the (single) metric Dimension
	DimensionValue *string

	// The namespace of the metric
	Namespace *string

	// The name of the metric
	MetricName *string

	// The average of the metric values that correspond to the data point.
	Average *float64

	// The maximum metric value for the data point.
	Maximum *float64

	// The minimum metric value for the data point.
	Minimum *float64

	// The number of metric values that contributed to the aggregate value of this
	// data point.
	SampleCount *float64

	// The sum of the metric values for the data point.
	Sum *float64

	// The time stamp used for the data point.
	Timestamp *time.Time

	// The standard unit for the data point.
	Unit *string

	Metadata map[string]string

	Region string
}

func getMonitoringStartDateForGranularity(granularity string) time.Time {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 90 days (We can fetch upto 90 days maximum)
		return time.Now().AddDate(0, 0, -90)
	case "HOURLY":
		// 60 days
		return time.Now().AddDate(0, 0, -60)
	}
	// else 5 days
	return time.Now().AddDate(0, 0, -5)
}

func getMonitoringPeriodForGranularity(granularity string) string {
	switch strings.ToUpper(granularity) {
	case "DAILY":
		// 24 hours
		return "1d"
	case "HOURLY":
		// 1 hour
		return "1h"
	}
	// else 5 minutes
	return "5m"
}

type MetricData struct {
	CompartmentId *string
	PointValue    *float64
	Timestamp     *time.Time
}

func listMonitoringMetricStatistics(ctx context.Context, d *plugin.QueryData, granularity string, namespace string, metricName string, dimensionName string, dimensionValue string, compartmentId string, region string) (*monitoring.SummarizeMetricsDataResponse, error) {
	plugin.Logger(ctx).Trace("listMonitoringMetricStatistics")

	// Create Session
	session, err := monitoringService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	/**
	DEFINE QUERY STRING
	metric[interval]{dimensionname="dimensionvalue"}.groupingfunction.statistic
	Query should be written with Metric query Language (MQL) https://docs.oracle.com/en-us/iaas/Content/Monitoring/Reference/mql.htm#Interval
	*/
	queryString := metricName + "[" + getMonitoringPeriodForGranularity(granularity) + "]" + "{" + dimensionName + " = \"" + dimensionValue + "\"}"
	queryStringavg := queryString + ".grouping().mean()"
	querystringMin := queryString + ".grouping().min()"
	querystringMax := queryString + ".grouping().max()"
	querystringSum := queryString + ".grouping().sum()"
	querystringCount := queryString + ".grouping().count()"

	// Set Inteval
	interval := getMonitoringPeriodForGranularity(granularity)
	metricDetails := monitoring.SummarizeMetricsDataDetails{
		Namespace:  &namespace,
		StartTime:  &common.SDKTime{Time: getMonitoringStartDateForGranularity(granularity)},
		EndTime:    &common.SDKTime{Time: time.Now()},
		Query:      &queryStringavg,
		Resolution: &interval,
	}

	requestParam := monitoring.SummarizeMetricsDataRequest{
		CompartmentId:               &compartmentId,
		SummarizeMetricsDataDetails: metricDetails,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Mean statistics
	avgStatistics, err := session.MonitoringClient.SummarizeMetricsData(ctx, requestParam)
	if err != nil {
		return nil, err
	}

	// Min statistics
	metricDetails.Query = &querystringMin
	requestParam.SummarizeMetricsDataDetails = metricDetails
	minStatistics, err := session.MonitoringClient.SummarizeMetricsData(ctx, requestParam)
	if err != nil {
		return nil, err
	}
	metricDetailsMin := filterMetricStatistic(minStatistics)

	// Max statistics
	metricDetails.Query = &querystringMax
	requestParam.SummarizeMetricsDataDetails = metricDetails
	maxStatistics, err := session.MonitoringClient.SummarizeMetricsData(ctx, requestParam)
	if err != nil {
		return nil, err
	}
	metricDetailsMax := filterMetricStatistic(maxStatistics)

	// Sum statistics
	metricDetails.Query = &querystringSum
	requestParam.SummarizeMetricsDataDetails = metricDetails
	sumStatistics, err := session.MonitoringClient.SummarizeMetricsData(ctx, requestParam)
	if err != nil {
		return nil, err
	}
	metricDetailsSum := filterMetricStatistic(sumStatistics)

	// Count statistics
	metricDetails.Query = &querystringCount
	requestParam.SummarizeMetricsDataDetails = metricDetails
	countStatistics, err := session.MonitoringClient.SummarizeMetricsData(ctx, requestParam)
	if err != nil {
		return nil, err
	}
	metricDetailsCount := filterMetricStatistic(countStatistics)

	for _, item := range avgStatistics.Items {
		for _, datapoint := range item.AggregatedDatapoints {
			d.StreamLeafListItem(ctx, &MonitoringMetricRow{
				CompartmentId:  item.CompartmentId,
				DimensionValue: &dimensionValue,
				DimensionName:  &dimensionName,
				Namespace:      &namespace,
				MetricName:     &metricName,
				Average:        datapoint.Value,
				Maximum:        getStatisticForColumnByTimestamp(datapoint.Timestamp.Time.UTC(), *item.CompartmentId, metricDetailsMax),
				Minimum:        getStatisticForColumnByTimestamp(datapoint.Timestamp.Time.UTC(), *item.CompartmentId, metricDetailsMin),
				Timestamp:      &datapoint.Timestamp.Time,
				SampleCount:    getStatisticForColumnByTimestamp(datapoint.Timestamp.Time.UTC(), *item.CompartmentId, metricDetailsCount),
				Sum:            getStatisticForColumnByTimestamp(datapoint.Timestamp.Time.UTC(), *item.CompartmentId, metricDetailsSum),
				Metadata:       item.Metadata,
				Region:         region,
			})
		}
	}

	return nil, err
}

func filterMetricStatistic(metricStatistic monitoring.SummarizeMetricsDataResponse) []MetricData {
	metricData := []MetricData{}
	for _, item := range metricStatistic.Items {
		for _, data := range item.AggregatedDatapoints {
			metricData = append(metricData, MetricData{
				CompartmentId: item.CompartmentId,
				PointValue:    data.Value,
				Timestamp:     &data.Timestamp.Time,
			})
		}
	}
	return metricData
}

func getStatisticForColumnByTimestamp(timestamp time.Time, compartmentId string, metricData []MetricData) *float64 {
	var value *float64
	for _, t := range metricData {
		if *t.Timestamp == timestamp && compartmentId == *t.CompartmentId {
			value = t.PointValue
			break
		}
	}

	return value
}
