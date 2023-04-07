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

func tableCoreVolumeBackupPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume_backup_policy",
		Description: "OCI Core Volume Backup Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreVolumeBackupPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumeBackupPolicies,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name for volume backup policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the volume backup policy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume backup policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns

			{
				Name:        "destination_region",
				Description: "The paired destination region for copying scheduled backups to.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "schedules",
				Description: "The collection of schedules that this policy will apply.",
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
				Transform:   transform.From(volumeBackupPolicyTags),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCoreVolumeBackupPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("core.listCoreVolumeBackupPolicies", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListVolumeBackupPoliciesRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListVolumeBackupPolicies(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumes := range response.Items {
			d.StreamListItem(ctx, volumes)

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

func getCoreVolumeBackupPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreVolumeBackupPolicy")
	logger := plugin.Logger(ctx)
	matrixRegion := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("core.getCoreVolumeBackupPolicy", "Compartment", compartment, "OCI_REGION", matrixRegion)

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
		return nil, err
	}

	request := core.GetVolumeBackupPolicyRequest{
		PolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BlockstorageClient.GetVolumeBackupPolicy(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.VolumeBackupPolicy, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Free-form tags
// 2. Defined Tags

func volumeBackupPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volumeBackupPolicy := d.HydrateItem.(core.VolumeBackupPolicy)

	var tags map[string]interface{}

	if volumeBackupPolicy.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volumeBackupPolicy.FreeformTags {
			tags[k] = v
		}
	}

	if volumeBackupPolicy.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeBackupPolicy.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
