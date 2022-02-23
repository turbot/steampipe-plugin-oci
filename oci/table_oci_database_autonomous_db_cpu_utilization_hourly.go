package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/database"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableOciDatabaseAutonomousDatabaseMetricCpuUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_autonomous_db_metric_cpu_utilization_hourly",
		Description: "OCI Autonomous Database Monitoring Metrics - CPU Utilization (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listAutonomousDatabases,
			Hydrate:       listAutonomousDatabaseMetricCpuUtilizationHourly,
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

func listAutonomousDatabaseMetricCpuUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	database := h.Item.(database.AutonomousDatabaseSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*database.Id))
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_autonomous_database", "CpuUtilization", "resourceId", *database.Id, *database.CompartmentId, region)
}
