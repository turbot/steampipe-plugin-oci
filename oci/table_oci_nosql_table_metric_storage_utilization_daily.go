package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/nosql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableOciNoSQLTableMetricStorageUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table_metric_storage_utilization_daily",
		Description: "OCI NoSQL Table Monitoring Metrics - Storage Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listNoSQLTables,
			Hydrate:       listNoSQLTableMetricStorageUtilizationDaily,
		},
		GetMatrixItem: BuildCompartementRegionList,
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

func listNoSQLTableMetricStorageUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	table := h.Item.(nosql.TableSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*table.Id))
	return listMonitoringMetricStatistics(ctx, d, "DAILY", "oci_nosql", "StorageGB", "tableName", *table.Name, *table.CompartmentId, region)
}
