package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_user",
		Description: "OCI Identity User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getUser,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The user's login for the Console.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "user_type",
				Description: "Type of the user. Value can be IDCS or IAM. Oracle Identity Cloud Service(IDCS) users authenticate through single sign-on and can be granted access to all services included in your account. IAM users can access Oracle Cloud Infrastructure services, but not all Cloud Platform services.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(userType),
			},
			{
				Name:        "time_created",
				Description: "Date and time the user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "The description assigned to the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The user's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_mfa_activated",
				Description: "The user's current state.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "can_use_console_password",
				Description: "Indicates if the user can log in to the console.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseConsolePassword"),
			},
			{
				Name:        "can_use_api_keys",
				Description: "Indicates if the user can use API keys.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseApiKeys"),
			},
			{
				Name:        "can_use_auth_tokens",
				Description: "Indicates if the user can use SWIFT passwords/auth tokens.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseAuthTokens"),
			},
			{
				Name:        "can_use_smtp_credentials",
				Description: "Indicates if the user can use SMTP passwords.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseSmtpCredentials"),
			},
			{
				Name:        "can_use_customer_secret_keys",
				Description: "Indicates if the user can use SigV4 symmetric keys.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseCustomerSecretKeys"),
			},
			{
				Name:        "can_use_o_auth2_client_credentials",
				Description: "Indicates if the user can use OAuth2 credentials and tokens.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Capabilities.CanUseOAuth2ClientCredentials"),
			},
			{
				Name:        "email",
				Description: "The email address you assign to the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email_verified",
				Description: "Whether the email address has been validated.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "identity_provider_id",
				Description: "The OCID of the `IdentityProvider` this user belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IdentityProviderId"),
			},
			{
				Name:        "inactive_status",
				Description: "Applicable only if the user's `lifecycleState` is INACTIVE. A 16-bit value showing the reason why the user is inactive. 0: SUSPENDED; 1: DISABLED; 2: BLOCKED (the user has exceeded the maximum number of failed login attempts for the Console)",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "external_identifier",
				Description: "Identifier of the user in the identity provider.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_groups",
				Description: "List of groups associated with the user.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getUserGroups,
				Transform:   transform.FromValue(),
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
				Transform:   transform.From(userTags),
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

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListUsersRequest{
		CompartmentId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListUsers(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, user := range response.Items {
			d.StreamListItem(ctx, user)
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

func getUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUser")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetUserRequest{
		UserId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.User, nil
}

func getUserGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(identity.User)
	plugin.Logger(ctx).Trace("getUserGroups")
	userGroups := []identity.UserGroupMembership{}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.ListUserGroupMembershipsRequest{
		CompartmentId: &session.TenancyID,
		UserId:        user.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListUserGroupMemberships(ctx, request)
		if err != nil {
			return nil, err
		}

		userGroups = append(userGroups, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return userGroups, nil
}

//// TRANSFORM FUNCTION

func userTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	user := d.HydrateItem.(identity.User)

	var tags map[string]interface{}

	if user.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range user.FreeformTags {
			tags[k] = v
		}
	}

	if user.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range user.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

func userType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	user := d.HydrateItem.(identity.User)

	if strings.Split(*user.Name, "/")[0] == "oracleidentitycloudservice" {
		return "IDCS", nil
	}

	return "IAM", nil
}
