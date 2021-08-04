package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreBootVolumeReplica(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_replica",
		Description: "OCI Core Boot Volume Replica",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreBootVolumeReplica,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreBootVolumeReplicas,
		},
		GetMatrixItem: BuildCompartementZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The boot volume replica's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of a boot volume replica.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "boot_volume_id",
				Description: "The OCID of the source boot volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of the boot volume replica.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the boot volume replica was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "image_id",
				Description: "The image OCID used to create the boot volume the replica is replicated from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ImageId"),
			},
			{
				Name:        "size_in_gbs",
				Description: "The size of the source boot volume, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInGBs"),
			},
			{
				Name:        "time_last_synced",
				Description: "The date and time the boot volume replica was last synced from the source boot volume.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeLastSynced.Time"),
			},
			{
				Name:        "total_data_transferred_in_gbs",
				Description: "The total size of the data transferred from the source boot volume to the boot volume replica, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("TotalDataTransferredInGBs"),
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bootVolumeReplicaTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(ociRegionName),
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCoreBootVolumeReplicas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreBootVolumeReplicas", "Compartment", compartment, "OCI_Zone", zone)

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListBootVolumeReplicasRequest{
		CompartmentId:      types.String(compartment),
		AvailabilityDomain: types.String(zone),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListBootVolumeReplicas(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, bootVolumeReplica := range response.Items {
			d.StreamListItem(ctx, bootVolumeReplica)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTION

func getCoreBootVolumeReplica(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreBootVolumeReplica", "Compartment", compartment, "OCI_Zone", zone)

	// Restrict the api call to only root compartment and one zone/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetBootVolumeReplicaRequest{
		BootVolumeReplicaId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BlockstorageClient.GetBootVolumeReplica(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.BootVolumeReplica, nil
}

//// TRANSFORM FUNCTION

func bootVolumeReplicaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volumeReplica := d.HydrateItem.(core.BootVolumeReplica)

	var tags map[string]interface{}

	if volumeReplica.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volumeReplica.FreeformTags {
			tags[k] = v
		}
	}

	if volumeReplica.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeReplica.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
