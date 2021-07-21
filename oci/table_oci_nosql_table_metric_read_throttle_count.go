package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/nosql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricReadThrottleCount(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_read_throttle_count",
		Description: "OCI NoSQL Table Monitoring Metrics - Read Throttle Count",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricReadThrottleCount,
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

func listNoSQLTableMetricReadThrottleCount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_nosql", "ReadThrottleCount", "tableName", *table.Name, *table.CompartmentId)
}
