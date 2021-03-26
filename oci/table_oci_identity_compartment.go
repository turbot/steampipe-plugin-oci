package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityCompartment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_compartment",
		Description: "OCI Identity Compartment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getCompartment,
		},
		List: &plugin.ListConfig{
			Hydrate: listCompartments,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name assigned to the compartment during creation",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the compartment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The compartment's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "description",
				Description: "The description you assign to the compartment.",
				Type:        proto.ColumnType_STRING,
			},

			// other columns
			{
				Name:        "inactive_status",
				Description: "The detailed status of INACTIVE lifecycleState",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_accessible",
				Description: "Indicates whether or not the compartment is accessible for the user making the request.",
				Type:        proto.ColumnType_BOOL,
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
				Transform:   transform.From(compartmentTags),
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
				Transform:   transform.FromCamel(),
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

func listCompartments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	rootRequest := identity.GetCompartmentRequest{
		CompartmentId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	responseRoot, err := session.IdentityClient.GetCompartment(ctx, rootRequest)
	if err != nil {
		return nil, err
	}

	if &responseRoot.Compartment != nil {
		d.StreamListItem(ctx, responseRoot.Compartment)
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListCompartmentsRequest{
		CompartmentId:          &session.TenancyID,
		CompartmentIdInSubtree: types.Bool(true),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListCompartments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, compartment := range response.Items {
			d.StreamListItem(ctx, compartment)
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

func getCompartment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCompartment")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetCompartmentRequest{
		CompartmentId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetCompartment(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Compartment, nil
}

//// TRANSFORM FUNCTION

func compartmentTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	compartment := d.HydrateItem.(identity.Compartment)

	var tags map[string]interface{}

	if compartment.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range compartment.FreeformTags {
			tags[k] = v
		}
	}

	if compartment.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range compartment.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}
	return tags, nil
}
