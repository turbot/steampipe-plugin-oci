package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityAuthenticationPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_authentication_policy",
		Description: "OCI Identity Authentication Policy",
		List: &plugin.ListConfig{
			Hydrate: listAuthenticationPolicy,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			// Password Policy
			{
				Name:        "is_lowercase_characters_required",
				Description: "At least one lower case character required.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PasswordPolicy.IsLowercaseCharactersRequired"),
			},
			{
				Name:        "is_numeric_characters_required",
				Description: "At least one numeric character required.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PasswordPolicy.IsNumericCharactersRequired"),
			},
			{
				Name:        "is_special_characters_required",
				Description: "At least one special character required.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PasswordPolicy.IsSpecialCharactersRequired"),
			},
			{
				Name:        "is_uppercase_characters_required",
				Description: "At least one uppercase character required.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PasswordPolicy.IsUppercaseCharactersRequired"),
			},
			{
				Name:        "is_username_containment_allowed",
				Description: "User name is allowed to be part of the password.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("PasswordPolicy.IsUsernameContainmentAllowed"),
			},
			{
				Name:        "minimum_password_length",
				Description: "Minimum password length required.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("PasswordPolicy.MinimumPasswordLength"),
			},

			// Network Policy
			{
				Name:        "network_source_ids",
				Description: "List of IP ranges from which users can sign in to the Console.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("NetworkPolicy.NetworkSourceIds"),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAuthenticationPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.GetAuthenticationPolicyRequest{
		CompartmentId:   &session.TenancyID,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetAuthenticationPolicy(ctx, request)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, response.AuthenticationPolicy)
	return nil, err
}
