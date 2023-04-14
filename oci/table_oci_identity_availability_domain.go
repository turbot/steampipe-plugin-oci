package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		Columns: commonColumnsForAllResource([]*plugin.Column{
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
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
		}),
	}
}

type availabilityDomainInfo struct {
	identity.AvailabilityDomain
	Region string
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
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.ListAvailabilityDomains(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, availabilityDomain := range response.Items {
		d.StreamListItem(ctx, availabilityDomainInfo{availabilityDomain, region})
	}

	return nil, nil
}
