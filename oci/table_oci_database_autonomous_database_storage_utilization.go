package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/database"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOciDatabaseAutonomousDatabaseMetricStorageUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_autonomous_database_metric_storage_utilization",
		Description: "OCI Autonomous Database Monitoring Metrics - Storage Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listAutonomousDatabases,
			Hydrate:       listAutonomousDatabaseMetricStorageUtilization,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The OCID of the Autonomous Database.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listAutonomousDatabaseMetricStorageUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(database.AutonomousDatabaseSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*database.Id))
	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_autonomous_database", "StorageUtilization", "resourceId", *database.Id, *database.CompartmentId, region)
}
