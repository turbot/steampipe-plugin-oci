package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityDomain(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_domain",
		Description: "OCI Identity Domain",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDomain,
		},
		List: &plugin.ListConfig{
			Hydrate: listDomains,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "url",
					Require: plugin.Optional,
				},
				{
					Name:    "home_region_url",
					Require: plugin.Optional,
				},
				{
					Name:    "type",
					Require: plugin.Optional,
				},
				{
					Name:    "license_type",
					Require: plugin.Optional,
				},
				{
					Name:    "is_hidden_on_login",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The mutable display name of the identity domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the identity domain.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "The identity domain description. You can have an empty description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the identity domain was created, in the format defined by RFC3339.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The domain's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Any additional details about the current state of the identity domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_hidden_on_login",
				Description: "Indicates whether the identity domain is hidden on the sign-in screen or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "url",
				Description: "Region-agnostic identity domain URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_region_url",
				Description: "Region-specific identity domain URL.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_region",
				Description: "The home region for the identity domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the identity domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_type",
				Description: "The license type of the identity domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "replica_regions",
				Description: "The regions where replicas of the identity domain exist.",
				Type:        proto.ColumnType_JSON,
			},
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
				Transform:   transform.From(domainTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listDomains(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	equalQuals := d.EqualsQuals

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListDomainsRequest{
		CompartmentId: &session.TenancyID,
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["diaplay_name"] != nil {
		name := equalQuals["diaplay_name"].GetStringValue()
		request.DisplayName = types.String(name)
	}
	if equalQuals["url"] != nil {
		url := equalQuals["url"].GetStringValue()
		request.Url = types.String(url)
	}
	if equalQuals["home_region_url"] != nil {
		homeRegionUrl := equalQuals["home_region_url"].GetStringValue()
		request.HomeRegionUrl = types.String(homeRegionUrl)
	}
	if equalQuals["type"] != nil {
		domainType := equalQuals["type"].GetStringValue()
		request.Type = types.String(domainType)
	}
	if equalQuals["license_type"] != nil {
		licenseType := equalQuals["license_type"].GetStringValue()
		request.Type = types.String(licenseType)
	}
	if equalQuals["is_hidden_on_login"] != nil {
		isLoginHidden := equalQuals["is_hidden_on_login"].GetBoolValue()
		request.IsHiddenOnLogin = types.Bool(isLoginHidden)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = identity.DomainLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListDomains(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, group := range response.Items {
			d.StreamListItem(ctx, group)

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

//// HYDRATE FUNCTIONS

func getDomain(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	id := d.EqualsQuals["id"].GetStringValue()

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetDomainRequest{
		DomainId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.GetDomain(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Domain, nil
}

//// TRANSFORM FUNCTION

func domainTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	switch domain := d.HydrateItem.(type) {
	case identity.Domain:
		return extractTags(domain.FreeformTags, domain.DefinedTags), nil
	case identity.DomainSummary:
		return extractTags(domain.FreeformTags, domain.DefinedTags), nil
	}
	return nil, nil
}
