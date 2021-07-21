package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/nosql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricStorageUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_storage_utilization",
		Description: "OCI NoSQL Table Monitoring Metrics - Storage Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricStorageUtilization,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "name",
					Description: "Immutable human-friendly table name.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listNoSQLTableMetricStorageUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_nosql", "StorageGB", "tableName", *table.Name, *table.CompartmentId)
}
