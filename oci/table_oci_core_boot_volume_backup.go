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

func tableCoreBootVolumeBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_backup",
		Description: "OCI Core Boot Volume Backup",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBootVolumeBackup,
		},
		List: &plugin.ListConfig{
			Hydrate: listBootVolumeBackups,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the boot volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name for the boot volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "boot_volume_id",
				Description: "The OCID of the boot volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of a boot volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expiration_time",
				Description: "The date and time the volume backup will expire and be automatically deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ExpirationTime.Time"),
			},
			{
				Name:        "image_id",
				Description: "The image OCID used to create the boot volume the backup is taken from.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the Key Management master encryption assigned to the boot volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the boot volume backup was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_request_received",
				Description: "The date and time the request to create the boot volume backup was received.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeRequestReceived.Time"),
			},
			{
				Name:        "type",
				Description: "The type of a volume backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size_in_gbs",
				Description: "The size of the boot volume, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInGBs"),
			},
			{
				Name:        "source_boot_volume_backup_id",
				Description: "The OCID of the source boot volume backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "source_type",
				Description: "Specifies whether the backup was created manually, or via scheduled backup policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "unique_size_in_gbs",
				Description: "The size used by the backup, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("UniqueSizeInGBs"),
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
				Description: "System tags for this resource.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bootVolumeBackupTags),
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

func listBootVolumeBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listBootVolumeBackups", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListBootVolumeBackupsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListBootVolumeBackups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, backup := range response.Items {
			d.StreamListItem(ctx, backup)
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

func getBootVolumeBackup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getBootVolumeBackup")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getBootVolumeBackup", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty boot volume backup id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetBootVolumeBackupRequest{
		BootVolumeBackupId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BlockstorageClient.GetBootVolumeBackup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.BootVolumeBackup, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func bootVolumeBackupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	bootVolumeBackup := d.HydrateItem.(core.BootVolumeBackup)

	var tags map[string]interface{}

	if bootVolumeBackup.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range bootVolumeBackup.FreeformTags {
			tags[k] = v
		}
	}

	if bootVolumeBackup.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range bootVolumeBackup.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	if bootVolumeBackup.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range bootVolumeBackup.SystemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
