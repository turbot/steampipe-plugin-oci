package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume",
		Description: "OCI Core Volume",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreVolume,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "volume_group_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the volume.",
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
				Description: "The current state of a volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "volume_group_id",
				Description: "The OCID of the source volume group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of the volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "auto_tuned_vpus_per_gb",
				Description: "The number of Volume Performance Units per GB that this volume is effectively tuned to when it's idle.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AutoTunedVpusPerGB"),
			},
			{
				Name:        "is_auto_tune_enabled",
				Description: "Specifies whether the auto-tune performance is enabled for this volume.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_hydrated",
				Description: "Specifies whether the cloned volume's data has finished copying from the source volume or backup.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the Key Management key which is the master encryption key for the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "size_in_gbs",
				Description: "The size of the volume in GBs.",
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
				Name:        "volume_backup_policy_id",
				Description: "The OCID of the volume backup policy that has been assigned to the volume.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVolumeBackupPolicyAssignment,
				Transform:   transform.FromField("PolicyId"),
			},
			{
				Name:        "volume_backup_policy_assignment_id",
				Description: "The OCID of the volume backup policy assignment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVolumeBackupPolicyAssignment,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "vpus_per_gb",
				Description: "The number of volume performance units (VPUs) that will be applied to this volume per GB,representing the Block Volume service's elastic performance options.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("VpusPerGB"),
			},

			// json fields
			{
				Name:        "source_details",
				Description: "The volume source, either an existing volume in the same availability domain or a volume backup.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(volumeTags),
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
			},
			{
				Name:        "compartment_id",
				Description: "ColumnDescriptionCompartment",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

type volumeInfo struct {
	core.Volume
	Region string
}

//// LIST FUNCTION

func listCoreVolumes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreVolumes", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildCoreVolumeFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListVolumes(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumes := range response.Items {
			d.StreamListItem(ctx, volumeInfo{volumes, region})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
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

func getCoreVolume(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreVolume")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreVolume", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty volume id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVolumeRequest{
		VolumeId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BlockstorageClient.GetVolume(ctx, request)
	if err != nil {
		return nil, err
	}

	return volumeInfo{response.Volume, region}, nil
}

func getVolumeBackupPolicyAssignment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVolumeBackupPolicyAssignment")
  region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)

  volumeId := h.Item.(volumeInfo).Id

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVolumeBackupPolicyAssetAssignmentRequest{
		AssetId: volumeId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BlockstorageClient.GetVolumeBackupPolicyAssetAssignment(ctx, request)
	if err != nil {
		plugin.Logger(ctx).Error("getVolumeBackupPolicyAssignment","err",err)
		return nil, err
	}

	return response.Items[0], nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func volumeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volume := d.HydrateItem.(volumeInfo).Volume

	var tags map[string]interface{}

	if volume.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volume.FreeformTags {
			tags[k] = v
		}
	}

	if volume.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volume.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	if volume.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volume.SystemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

// Build additional filters
func buildCoreVolumeFilters(equalQuals plugin.KeyColumnEqualsQualMap) core.ListVolumesRequest {
	request := core.ListVolumesRequest{}

	if equalQuals["availability_domain"] != nil {
		request.AvailabilityDomain = types.String(equalQuals["availability_domain"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = core.VolumeLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}
	if equalQuals["volume_group_id"] != nil {
		request.VolumeGroupId = types.String(equalQuals["volume_group_id"].GetStringValue())
	}

	return request
}
