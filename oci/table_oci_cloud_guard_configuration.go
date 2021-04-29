package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v36/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudGuardConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_configuration",
		Description: "OCI Cloud Guard Configuration",
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardConfigurations,
		},
		Columns: []*plugin.Column{
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudGuardConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetConfigurationRequest{
		CompartmentId: types.String(session.TenancyID),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.CloudGuardClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, response.Configuration)

	return nil, err
}
