package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOciMySQLDBSystemMetricCpuUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_db_system_metric_cpu_utilization",
		Description: "OCI MySQL DB System Monitoring Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listMySQLDBSystems,
			Hydrate:       listMySQLDBSystemMetricCpuUtilization,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The OCID of the DB System.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listMySQLDBSystemMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbSystem := h.Item.(mysql.DbSystemSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*dbSystem.Id))

	if dbSystem.LifecycleState == "DELETING" || dbSystem.LifecycleState == "DELETED" {
		return nil, nil
	}

	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_mysql_database", "CPUUtilization", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
}
