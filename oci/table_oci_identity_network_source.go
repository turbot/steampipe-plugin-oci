package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityNetworkSource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_network_source",
		Description: "OCI Identity Network Source",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getIdentityNetworkSource,
		},
		List: &plugin.ListConfig{
			Hydrate: listIdentityNetworkSources,
		},
		Columns: []*plugin.Column{
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
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
		},
	}
}

//// LIST FUNCTION

func listIdentityNetworkSources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListNetworkSourcesRequest{
		CompartmentId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListNetworkSources(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, networkSources := range response.Items {
			d.StreamListItem(ctx, networkSources)
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
	plugin.Logger(ctx).Trace("getIdentityNetworkSource")

	var id string
	if h.Item != nil {
		id = *h.Item.(identity.NetworkSourcesSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty network source id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetNetworkSourceRequest{
		NetworkSourceId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
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
