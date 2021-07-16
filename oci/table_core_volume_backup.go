package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVolumeBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume_backup",
		Description: "OCI Core Volume Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getVolumeBackup,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumeBackups,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name for the volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "volume_id",
				Description: "The OCID of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "expiration_time",
				Description: "The date and time the volume backup will expire and be automatically deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpirationTime.Time"),
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the Key Management key which is the master encryption key for the volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of a volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size_in_gbs",
				Description: "The size of the volume, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInGBs"),
			},
			{
				Name:        "size_in_mbs",
				Description: "The size of the volume in MBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInMBs"),
			},
			{
				Name:        "source_type",
				Description: "Specifies whether the backup was created manually, or via scheduled backup policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_volume_backup_id",
				Description: "The OCID of the source volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume backup was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "type",
				Description: "The type of a volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_request_received",
				Description: "The date and time the request to create the volume backup was received.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeRequestReceived.Time"),
			},
			{
				Name:        "unique_size_in_gbs",
				Description: "The size used by the backup, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("UniqueSizeInGBs"),
			},
			{
				Name:        "unique_size_in_mbs",
				Description: "The size used by the backup, in MBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("UniqueSizeInMbs"),
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
			{
				Name:        "system_tags",
				Description: "System tags to volume by the service.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(volumeBackupTags),
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

func listCoreVolumeBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreVolumeBackups", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListVolumeBackupsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListVolumeBackups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeBackups := range response.Items {
			d.StreamListItem(ctx, volumeBackups)
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

func getVolumeBackup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVolumeBackup")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getVolumeBackup", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty volume backup id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVolumeBackupRequest{
		VolumeBackupId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BlockstorageClient.GetVolumeBackup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.VolumeBackup, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func volumeBackupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volumeBackup := d.HydrateItem.(core.VolumeBackup)

	var tags map[string]interface{}

	if volumeBackup.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volumeBackup.FreeformTags {
			tags[k] = v
		}
	}

	if volumeBackup.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeBackup.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	if volumeBackup.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeBackup.SystemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
