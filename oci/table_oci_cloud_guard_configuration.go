package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/cloudguard"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudGuardConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_configuration",
		Description: "OCI Cloud Guard Configuration",
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardConfigurations,
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "reporting_region",
				Description: "The reporting region value.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Status of Cloud Guard Tenant.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "self_manage_resources",
				Description: "Identifies if Oracle managed resources were created by customers.",
				Type:        proto.ColumnType_BOOL,
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

func listCloudGuardConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getCloudGuardConfigurationCached := plugin.HydrateFunc(getCloudGuardConfiguration).WithCache()
	configuration, err := getCloudGuardConfigurationCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, configuration.(cloudguard.Configuration))

	return nil, nil
}
