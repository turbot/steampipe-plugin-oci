package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
					Name:        "volume_id",
					Description: "The OCID of the Boot Volume.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			}),
	}
}

func listCoreBootVolumeMetricWriteOpsDaily(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// We can get the metric details of volume, only which are attach with the instance or boot volume
	volume := h.Item.(core.BootVolume)
	return listMonitoringMetricStastics(ctx, d, "DAILY", "oci_blockstore", "VolumeWriteOps", "resourceId", *volume.Id, *volume.CompartmentId)
}
