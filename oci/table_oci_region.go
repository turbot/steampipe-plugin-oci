package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityRegion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_region",
		Description: "OCI Region",
		List: &plugin.ListConfig{
			Hydrate: listRegions,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the region.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "key",
				Description: "The key of the region.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Key"),
			},
			{
				Name:        "is_home_region",
				Description: "Indicates if the region is the home region or not.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Transform:   transform.FromField("IsHomeRegion"),
			},
			{
				Name:        "status",
				Description: "The region subscription status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status"),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

type ociRegion struct {
	identity.Region
	identity.RegionSubscription
}

//// LIST FUNCTION

func listRegions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// List all the regions for the tenant
	regions, err := session.IdentityClient.ListRegions(ctx)
	if err != nil {
		return nil, err
	}

	request := identity.ListRegionSubscriptionsRequest{
		TenancyId: &session.TenancyID,
	}

	// List all the subscribed regions for the tenant
	subscribedRegions, err := session.IdentityClient.ListRegionSubscriptions(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, region := range regions.Items {
		isSubscribed := false
		for _, subscribedRegion := range subscribedRegions.Items {
			if *region.Name == *subscribedRegion.RegionName {
				d.StreamListItem(ctx, ociRegion{region, subscribedRegion})
				isSubscribed = true
				break
			}
		}
		if isSubscribed {
			continue
		}
		d.StreamListItem(ctx, ociRegion{region, identity.RegionSubscription{}})
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getTenantId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTenantId")
	cacheKey := "getTenantId"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// cache tenant id for the session
	d.ConnectionManager.Cache.Set(cacheKey, session.TenancyID)

	return session.TenancyID, nil
}
