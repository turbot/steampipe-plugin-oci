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
		GetMatrixItem: BuildCompartementRegionList,
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
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ReportingRegion"),
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

func listCloudGuardConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listCloudGuardConfigurations", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := cloudGuardService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetConfigurationRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}
	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudGuardClient.GetConfiguration(ctx, request)
		if err != nil {
			return nil, err
		}
		d.StreamListItem(ctx, response.Configuration)
		pagesLeft = false
	}

	return nil, err
}
