package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableOciMySQLDBSystemMetricConnections(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_db_system_metric_connections",
		Description: "OCI MySQL DB System Monitoring Metrics - Current/Active Connections",
		List: &plugin.ListConfig{
			ParentHydrate: listMySQLDBSystems,
			Hydrate:       listMySQLDBSystemMetricConnections,
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

func listMySQLDBSystemMetricConnections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	dbSystem := h.Item.(mysql.DbSystemSummary)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*dbSystem.Id))

	if dbSystem.LifecycleState == "DELETING" || dbSystem.LifecycleState == "DELETED" {
		return nil, nil
	}

	_, err := listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_mysql_database", "ActiveConnections", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
	if err != nil {
		return nil, err
	}

	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_mysql_database", "CurrentConnections", "resourceId", *dbSystem.Id, *dbSystem.CompartmentId, region)
}
