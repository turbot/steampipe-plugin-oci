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
func tableOciCoreBootVolumeMetricReadOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_metric_read_ops_daily",
		Description: "OCI Core Boot Volume Monitoring Metrics - Read Ops (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listBootVolumes,
			Hydrate:       listCoreBootVolumeMetricReadOpsDaily,
		},
		GetMatrixItemFunc: BuildCompartmentZonalList,
		Columns: MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The OCID of the boot volume.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listCoreBootVolumeMetricReadOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// We can get the metric details of volume, only which are attach with the instance or boot volume
	volume := h.Item.(core.BootVolume)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*volume.Id))
	return listMonitoringMetricStatistics(ctx, d, "DAILY", "oci_blockstore", "VolumeReadOps", "resourceId", *volume.Id, *volume.CompartmentId, region)
}
