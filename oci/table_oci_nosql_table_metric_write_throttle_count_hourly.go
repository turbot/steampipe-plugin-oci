package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/nosql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricWriteThrottleCountHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_write_throttle_count_hourly",
		Description: "OCI NoSQL Table Monitoring Metrics - Write Throttle Count (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricWriteThrottleCountHourly,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "table_name",
					Description: "The NoSQL Table Name.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listNoSQLTableMetricWriteThrottleCountHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_nosql", "WriteThrottleCount", "tableName", *table.Name, *table.CompartmentId)
}
