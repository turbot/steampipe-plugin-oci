package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/nosql"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricStorageUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_storage_utilization_hourly",
		Description: "OCI NoSQL Table Monitoring Metrics - Storage Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricStorageUtilizationHourly,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource(MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "The name of the NoSQL table.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listNoSQLTableMetricStorageUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*table.Id))
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_nosql", "StorageGB", "tableName", *table.Name, *table.CompartmentId, region)
}
