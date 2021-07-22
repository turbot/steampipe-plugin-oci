package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION
func tableOciCoreInstanceMetricCpuUtilizationHourly(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_instance_metric_cpu_utilization_hourly",
		Description: "OCI Core Instance Monitoring Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listCoreInstances,
			Hydrate:       listCoreInstanceMetricCpuUtilizationHourly,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The OCID of the instance.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listCoreInstanceMetricCpuUtilizationHourly(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(core.Instance)
	return listMonitoringMetricStatistics(ctx, d, "HOURLY", "oci_computeagent", "CpuUtilization", "resourceId", *instance.Id, *instance.CompartmentId, "")
}