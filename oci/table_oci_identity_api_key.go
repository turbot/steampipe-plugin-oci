package oci

import (
	"context"
	"strconv"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityApiKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_api_key",
		Description: "OCI Identity API Key",
		List: &plugin.ListConfig{
			ParentHydrate: listUsers,
			Hydrate:       listIdentityApiKeys,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "user_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "key_id",
				Description: "An Oracle-assigned identifier for the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "key_value",
				Description: "The key's value.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "The OCID of the user the key belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_name",
				Description: "The name of the user the key belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The API key's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fingerprint",
				Description: "The key's fingerprint.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the `ApiKey` object was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Transform:   transform.FromField("KeyId"),
			},

			// Standard OCI columns
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type apiKeyInfo struct {
	identity.ApiKey
	UserName string
}

//// LIST FUNCTION

func listIdentityApiKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	// The OCID of the User.
	request := identity.ListApiKeysRequest{
		UserId: user.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// List user's API key
	item, err := session.IdentityClient.ListApiKeys(ctx, request)
	if err != nil {
		if ociErr, ok := err.(oci_common.ServiceError); ok {
			if helpers.StringSliceContains([]string{"404"}, strconv.Itoa(ociErr.GetHTTPStatusCode())) {
				return nil, nil
			}
		}
		return nil, err
	}

	for _, apiKey := range item.Items {
		d.StreamListItem(ctx, apiKeyInfo{apiKey, *user.Name})
	}

	return nil, nil
}
