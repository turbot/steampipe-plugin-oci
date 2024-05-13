package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityAuthToken(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_auth_token",
		Description: "OCI Identity Auth Token",
		List: &plugin.ListConfig{
			ParentHydrate: listUsers,
			Hydrate:       listIdentityAuthTokens,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "user_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the auth token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_id",
				Description: "The OCID of the user the auth token belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_name",
				Description: "The name of the user the auth token belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "token",
				Description: "The auth token. The value is available only in the response for `CreateAuthToken`, and not for `ListAuthTokens` or `UpdateAuthToken`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The token's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the `AuthToken` object was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_expires",
				Description: "Date and time when this auth token will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeExpires.Time"),
			},

			// other columns
			{
				Name:        "description",
				Description: "The description you assign to the auth token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inactive_status",
				Description: "The detailed status of INACTIVE lifecycleState.",
				Type:        proto.ColumnType_INT,
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type authTokenInfo struct {
	identity.AuthToken
	UserName string
}

//// LIST FUNCTION

func listIdentityAuthTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(identity.User)

	// Return nil, if given user_id doesn't match
	equalQuals := d.EqualsQuals
	if equalQuals["user_id"] != nil && equalQuals["user_id"].GetStringValue() != *user.Id {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListAuthTokensRequest{
		UserId: user.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// List user auth tokens
	item, err := session.IdentityClient.ListAuthTokens(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, authToken := range item.Items {
		d.StreamLeafListItem(ctx, authTokenInfo{authToken, *user.Name})
	}

	return nil, nil
}
