package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVolumeDefaultBackupPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume_default_backup_policy",
		Description: "OCI Core Volume Default Backup Policy",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreVolumeDefaultBackupPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumeDefaultBackupPolicies,
		},
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

			// json fields
			{
				Name:        "schedules",
				Description: "The collection of schedules that this policy will apply.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCoreVolumeDefaultBackupPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	// Create Session
	session, err := coreBlockStorageService(ctx, d, "")
	if err != nil {
		logger.Error("oci_core_volume_default_backup_policy.listCoreVolumeDefaultBackupPolicies", "connection_error", err)
		return nil, err
	}

	request := core.ListVolumeBackupPoliciesRequest{
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
			logger.Error("oci_core_volume_default_backup_policy.listCoreVolumeDefaultBackupPolicies", "api_error", err)
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

func getCoreVolumeDefaultBackupPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	id := d.EqualsQuals["id"].GetStringValue()

	// handle empty volume backup policy id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, "")
	if err != nil {
		logger.Error("oci_core_volume_default_backup_policy.getCoreVolumeDefaultBackupPolicy", "connection_error", err)
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
		logger.Error("oci_core_volume_default_backup_policy.getCoreVolumeDefaultBackupPolicy", "api_error", err)
		return nil, err
	}

	return response.VolumeBackupPolicy, nil
}
