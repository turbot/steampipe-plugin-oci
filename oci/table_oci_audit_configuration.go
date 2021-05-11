package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v36/audit"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAuditConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_audit_configuration",
		Description: "OCI Audit Configuration",
		List: &plugin.ListConfig{
			Hydrate: listAuditConfiguration,
		},
		Columns: []*plugin.Column{
			{
				Name:        "retention_period_days",
				Description: "The retention period setting, specified in days. The minimum is 90, the maximum 365.",
				Type:        proto.ColumnType_INT,
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

func listAuditConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := auditService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := audit.GetConfigurationRequest{
		CompartmentId: types.String(session.TenancyID),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.AuditClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, response.Configuration)

	return nil, err
}
