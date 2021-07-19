package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityAvailabilityDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_availability_domain",
		Description: "OCI Identity Availability Domain",
		List: &plugin.ListConfig{
			ParentHydrate: listRegions,
			Hydrate:       lisAvailabilityDomains,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the Availability Domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the Availability Domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard OCI columns
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
		},
	}
}

//// LIST FUNCTION

func lisAvailabilityDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("lisAvailabilityDomains")

	region := *h.Item.(ociRegion).Name
	status := h.Item.(ociRegion).Status

	// Check if the region is subscribed region
	if status != "READY" {
		return nil, nil
	}

	// Create Session
	session, err := identityServiceRegional(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: &session.TenancyID,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.ListAvailabilityDomains(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, availabilityDomain := range response.Items {
		d.StreamListItem(ctx, availabilityDomain)
	}

	return nil, nil
}
