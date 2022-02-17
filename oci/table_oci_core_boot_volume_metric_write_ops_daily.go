package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION
func tableOciCoreBootVolumeMetricWriteOpsDaily(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_metric_write_ops_daily",
		Description: "OCI Core Boot Volume Monitoring Metrics - Write Ops (Daily)",
		List: &plugin.ListConfig{
			ParentHydrate: listBootVolumes,
			Hydrate:       listCoreBootVolumeMetricWriteOpsDaily,
		},
		GetMatrixItem: BuildCompartementZonalList,
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

func listCoreBootVolumeMetricWriteOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// We can get the metric details of volume, only which are attach with the instance or boot volume
	volume := h.Item.(core.BootVolume)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*volume.Id))
	return listMonitoringMetricStatistics(ctx, d, "DAILY", "oci_blockstore", "VolumeWriteOps", "resourceId", *volume.Id, *volume.CompartmentId, region)
}
