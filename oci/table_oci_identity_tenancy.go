package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/audit"
	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityTenancy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_tenancy",
		Description: "OCI Identity Tenancy",
		List: &plugin.ListConfig{
			Hydrate: listIdentityTenancies,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:              getRetentionPeriod,
				ShouldIgnoreError: isNotFoundError([]string{"404"}),
			},
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the tenancy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the tenancy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "retention_period_days",
				Description: "The retention period setting, specified in days.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getRetentionPeriod,
			},
			{
				Name:        "description",
				Description: "The description of the tenancy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_region_key",
				Description: "The region key for the tenancy's home region.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "upi_idcs_compatibility_layer_endpoint",
				Description: "Url which refers to the UPI IDCS compatibility layer endpoint configured for this Tenant's home region.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UpiIdcsCompatibilityLayerEndpoint"),
			},

			// tags
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(tenantTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard OCI columns
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
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

func listIdentityTenancies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.GetTenancyRequest{
		TenancyId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.GetTenancy(ctx, request)
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, response.Tenancy)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getRetentionPeriod(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getRetentionPeriod")

	// Create Session
	session, err := auditService(ctx, d)
	if err != nil {
		return nil, err
	}

	var compartmentID string
	if h.Item != nil {
		compartmentID = *h.Item.(identity.Tenancy).Id
	} else {
		compartmentID = session.TenancyID
	}

	request := audit.GetConfigurationRequest{
		CompartmentId: &compartmentID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AuditClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Configuration, nil
}

//// TRANSFORM FUNCTION

func tenantTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	tenancy := d.HydrateItem.(identity.Tenancy)

	var tags map[string]interface{}

	if tenancy.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range tenancy.FreeformTags {
			tags[k] = v
		}
	}

	if tenancy.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range tenancy.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
