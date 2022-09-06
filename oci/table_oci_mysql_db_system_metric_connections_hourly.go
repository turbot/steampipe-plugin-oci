package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableOciMySQLDBSystemMetricConnectionsHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_db_system_metric_connections_hourly",
		Description: "OCI MySQL DB System Monitoring Metrics - Current/Active Connections (Hourly)",
		List: &plugin.ListConfig{
			ParentHydrate: listMySQLDBSystems,
			Hydrate:       listMySQLDBSystemMetricConnectionsHourly,
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

func listMySQLDBSystemMetricConnectionsHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbSystem := h.Item.(mysql.DbSystemSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*dbSystem.Id))

	if dbSystem.LifecycleState == "DELETING" || dbSystem.LifecycleState == "DELETED" {
		return nil, nil
	}

	_, err := listMonitoringMetricStatistics(ctx, d, "Hourly", "oci_mysql_database", "ActiveConnections", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
	if err != nil {
		return nil, err
	}

	return listMonitoringMetricStatistics(ctx, d, "Hourly", "oci_mysql_database", "CurrentConnections", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
}
