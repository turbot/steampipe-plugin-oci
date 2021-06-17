package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityApiKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_api_key",
		Description: "OCI Identity API Key",
		List: &plugin.ListConfig{
			ParentHydrate: listUsers,
			Hydrate:       listIdentityApiKeys,
		},
		Columns: []*plugin.Column{
			{
				Name:        "key_id",
				Description: "An Oracle-assigned identifier for the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "key_value",
				Description: "The key's value..",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_id",
				Description: "The OCID of the user the password belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_name",
				Description: "The name of the user the password belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The API key's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fingerprint",
				Description: " The key's fingerprint.",
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

type apiKeyInfo struct {
	identity.ApiKey
	UserName string
}

//// LIST FUNCTION

func listIdentityApiKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	user := h.Item.(identity.User)

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListApiKeysRequest{
		UserId: user.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// List user's API key
	item, err := session.IdentityClient.ListApiKeys(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, apiKey := range item.Items {
		d.StreamLeafListItem(ctx, apiKeyInfo{apiKey, *user.Name})
	}

	return nil, nil
}
