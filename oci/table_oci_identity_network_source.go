package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityNetworkSource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_network_source",
		Description: "OCI Identity Network Source",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIdentityNetworkSource,
		},
		List: &plugin.ListConfig{
			Hydrate: listIdentityNetworkSources,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartmentList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name you assign to the network source during creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the network source.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The network source object's current state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIdentityNetworkSource,
			},
			{
				Name:        "time_created",
				Description: "Date and time the etwork source was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "The description you assign to the network source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inactive_status",
				Description: "The detailed status of INACTIVE lifecycleState.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getIdentityNetworkSource,
			},

			// json fields
			{
				Name:        "public_source_list",
				Description: "A list of allowed public IP addresses and CIDR ranges.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "services",
				Description: "A list of services allowed to make on-behalf-of requests.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "virtual_source_list",
				Description: "A list of allowed VCN OCID and IP range pairs.",
				Type:        proto.ColumnType_JSON,
			},

			// tags
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
				Hydrate:     getIdentityNetworkSource,
				Transform:   transform.From(networkSourceTags),
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
				Transform:   transform.FromField("CompartmentId"),
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

//// LIST FUNCTION

func listIdentityNetworkSources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	compartment := d.EqualsQualString(matrixKeyCompartment)
	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListNetworkSourcesRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["name"] != nil {
		name := equalQuals["name"].GetStringValue()
		request.Name = types.String(name)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = identity.NetworkSourcesLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListNetworkSources(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, networkSources := range response.Items {
			d.StreamListItem(ctx, networkSources)

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

func getIdentityNetworkSource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	compartment := d.EqualsQualString(matrixKeyCompartment)

	var id string
	if h.Item != nil {
		id = *h.Item.(identity.NetworkSourcesSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// Restrict the api call to only root compartment
	// Handle empty dynamic group id in get call
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || id == "" {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetNetworkSourceRequest{
		NetworkSourceId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.GetNetworkSource(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.NetworkSources, nil
}

//// TRANSFORM FUNCTION

func networkSourceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	networkSource := d.HydrateItem.(identity.NetworkSources)

	var tags map[string]interface{}

	if networkSource.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range networkSource.FreeformTags {
			tags[k] = v
		}
	}

	if networkSource.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range networkSource.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
