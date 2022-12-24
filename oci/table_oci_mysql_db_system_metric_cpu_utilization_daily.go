package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/mysql"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableOciMySQLDBSystemMetricCpuUtilizationDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_db_system_metric_cpu_utilization_daily",
		Description: "OCI MySQL DB System Monitoring Metrics - CPU Utilization (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listMySQLDBSystems,
			Hydrate:       listMySQLDBSystemMetricCpuUtilizationDaily,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
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

func listMySQLDBSystemMetricCpuUtilizationDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbSystem := h.Item.(mysql.DbSystemSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*dbSystem.Id))

	if dbSystem.LifecycleState == "DELETING" || dbSystem.LifecycleState == "DELETED" {
		return nil, nil
	}

	return listMonitoringMetricStatistics(ctx, d, "DAILY", "oci_mysql_database", "CPUUtilization", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
}
