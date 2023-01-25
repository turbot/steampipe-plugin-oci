package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/nosql"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricReadThrottleCountHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_read_throttle_count_hourly",
		Description: "OCI NoSQL Table Monitoring Metrics - Read Throttle Count (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricReadThrottleCountHourly,
		},
		GetMatrixItemFunc: BuildCompartmentRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The name of the NoSQL table.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listNoSQLTableMetricReadThrottleCountHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*table.Id))
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_nosql", "ReadThrottleCount", "tableName", *table.Name, *table.CompartmentId, region)
}
