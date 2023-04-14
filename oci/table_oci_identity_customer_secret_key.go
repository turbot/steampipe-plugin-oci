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

func tableIdentityCustomerSecretKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_customer_secret_key",
		Description: "OCI Identity Customer Secret Key",
		List: &plugin.ListConfig{
			ParentHydrate: listUsers,
			Hydrate:       listIdentityCustomerSecretKeys,
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the secret key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "The displayName you assign to the secret key.",
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
				Description: "The secret key's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the CustomerSecretKey object was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_expires",
				Description: "Date and time when this password will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeExpires.Time"),
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
				Transform:   transform.FromField("DisplayName"),
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

type customerSecretKeyInfo struct {
	identity.CustomerSecretKeySummary
	UserName string
}

//// LIST FUNCTION

func listIdentityCustomerSecretKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	user := h.Item.(identity.User)

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListCustomerSecretKeysRequest{
		UserId: user.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// List user's customer secret key
	item, err := session.IdentityClient.ListCustomerSecretKeys(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, secretKey := range item.Items {
		d.StreamLeafListItem(ctx, customerSecretKeyInfo{secretKey, *user.Name})
	}

	return nil, nil
}
