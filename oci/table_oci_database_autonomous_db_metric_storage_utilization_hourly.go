package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableOciDatabaseAutonomousDatabaseMetricStorageUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_autonomous_db_metric_storage_utilization_hourly",
		Description: "OCI Autonomous Database Monitoring Metrics - Storage Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listAutonomousDatabases,
			Hydrate:       listAutonomousDatabaseMetricStorageUtilizationHourly,
		},
		GetMatrixItemFunc: BuildCompartmentRegionList,
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

func listAutonomousDatabaseMetricStorageUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(database.AutonomousDatabaseSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*database.Id))
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_autonomous_database", "StorageUtilization", "resourceId", *database.Id, *database.CompartmentId, region)
}
