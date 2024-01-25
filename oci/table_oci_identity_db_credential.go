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

func tableIdentityDBCredential(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_db_credential",
		Description: "OCI Identity DB Credential",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDomain,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listUsers,
			Hydrate:       listIdentityDBCredentials,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "user_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the DB credential.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_id",
				Description: "The OCID of the user the DB credential belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "The description you assign to the DB credential. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the DbCredential object was created, in the format defined by RFC3339.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_expires",
				Description: "Date and time when this credential will expire, in the format defined by RFC3339.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "lifecycle_state",
				Description: "The credential's current state.",
				Type:        proto.ColumnType_STRING,
			},
			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},

			// Standard OCI columns
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

func listIdentityDBCredentials(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	equalQuals := d.EqualsQuals
	user := h.Item.(identity.User)

	// Minimize API call with given User ID.
	if d.EqualsQualString("user_id") != "" && d.EqualsQualString("user_id") != *user.Id {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListDbCredentialsRequest{
		UserId: user.Id,
		Limit:  types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := d.EqualsQualString("lifecycle_state")
		request.LifecycleState = identity.DbCredentialLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListDbCredentials(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, credential := range response.Items {
			d.StreamListItem(ctx, credential)

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
