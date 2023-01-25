package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION
func tableOciCoreInstanceMetricCpuUtilization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_instance_metric_cpu_utilization",
		Description: "OCI Core Instance Monitoring Metrics - CPU Utilization",
		List: &plugin.ListConfig{
			ParentHydrate: listCoreInstances,
			Hydrate:       listCoreInstanceMetricCpuUtilization,
		},
		GetMatrixItemFunc: BuildCompartmentRegionList,
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

func listCoreInstanceMetricCpuUtilization(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	instance := h.Item.(core.Instance)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*instance.Id))
	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_computeagent", "CpuUtilization", "resourceId", *instance.Id, *instance.CompartmentId, region)
}
