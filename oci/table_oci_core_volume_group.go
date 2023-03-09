package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVolumeGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume_group",
		Description: "OCI Core Volume Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreVolumeGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumeGroups,
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
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID for the volume group.",
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
				Description: "The current state of a volume group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of the volume group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "is_hydrated",
				Description: "Specifies whether the cloned volume's data has finished copying from the source volume group or backup.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "size_in_gbs",
				Description: "The aggregate size of the volume group in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInGBs"),
			},
			{
				Name:        "size_in_mbs",
				Description: "The aggregate size of the volume group in MBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInMBs"),
			},

			// json fields
			{
				Name:        "source_details",
				Description: "The volume group source, either an existing volume group in the same availability domain or a volume group backup.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "volume_ids",
				Description: "OCIDs for the volumes in this volume group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "volume_group_replicas",
				Description: "The list of volume group replicas of this volume group.",
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(volumeGroupTags),
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

//// LIST FUNCTION

func listCoreVolumeGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("oci_core_volume_group.listCoreVolumeGroups", "session_error", err)
		return nil, err
	}

	// Build request parameters
	request := buildCoreVolumeGroupFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListVolumeGroups(ctx, request)
		if err != nil {
			plugin.Logger(ctx).Error("oci_core_volume_group.listCoreVolumeGroups", "api_error", err)
			return nil, err
		}

		for _, volumeGroup := range response.Items {
			d.StreamListItem(ctx, volumeGroup)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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

func getCoreVolumeGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	matrixRegion := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	id := d.EqualsQuals["id"].GetStringValue()

	// For the us-phoenix-1 and us-ashburn-1 regions, `phx` and `iad` are returned by ListInstances api, respectively.
	// For all other regions, the full region name is returned.
	region := common.StringToRegion(types.SafeString(strings.Split(id, ".")[3]))

	// handle empty id and region check in get call
	if id == "" || region != common.StringToRegion(matrixRegion) {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, matrixRegion)
	if err != nil {
		plugin.Logger(ctx).Error("oci_core_volume_group.getCoreVolumeGroup", "session_error", err)
		return nil, err
	}

	request := core.GetVolumeGroupRequest{
		VolumeGroupId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BlockstorageClient.GetVolumeGroup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.VolumeGroup, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func volumeGroupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volumeGroup := d.HydrateItem.(core.VolumeGroup)

	var tags map[string]interface{}

	if volumeGroup.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volumeGroup.FreeformTags {
			tags[k] = v
		}
	}

	if volumeGroup.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeGroup.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

// Build additional filters
func buildCoreVolumeGroupFilters(equalQuals plugin.KeyColumnEqualsQualMap) core.ListVolumeGroupsRequest {
	request := core.ListVolumeGroupsRequest{}

	if equalQuals["availability_domain"] != nil {
		request.AvailabilityDomain = types.String(equalQuals["availability_domain"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = core.VolumeGroupLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
