package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION
func tableOciCoreBootVolumeMetricReadOps(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_metric_read_ops",
		Description: "OCI Core Boot Volume Monitoring Metrics - Read Ops",
		List: &plugin.ListConfig{
			ParentHydrate: listBootVolumes,
			Hydrate:       listCoreBootVolumeMetricReadOps,
		},
		GetMatrixItemFunc: BuildCompartementZonalList,
		Columns: commonColumnsForAllResource(MonitoringMetricColumns(
			[]*plugin.Column{
				{
					Name:        "id",
					Description: "The OCID of the boot volume.",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromField("DimensionValue"),
				},
			})),
	}
}

func listCoreBootVolumeMetricReadOps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// We can get the metric details of volume, only which are attach with the instance or boot volume
	volume := h.Item.(core.BootVolume)
	region := fmt.Sprintf("%v", ociRegionNameFromId(*volume.Id))
	return listMonitoringMetricStatistics(ctx, d, "5_MIN", "oci_blockstore", "VolumeReadOps", "resourceId", *volume.Id, *volume.CompartmentId, region)
}
