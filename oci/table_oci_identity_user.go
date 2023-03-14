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

func tableIdentityUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_user",
		Description: "OCI Identity User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "external_identifier",
					Require: plugin.Optional,
				},
				{
					Name:    "identity_provider_id",
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
	equalQuals := d.EqualsQuals

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildUserGroupFilters(equalQuals)
	request.CompartmentId = &session.TenancyID
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListUsers(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, user := range response.Items {
			d.StreamListItem(ctx, user)

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

func getUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUser")

	id := d.EqualsQuals["id"].GetStringValue()

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetUserRequest{
		UserId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
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
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
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

// Build additional filters
func buildUserGroupFilters(equalQuals plugin.KeyColumnEqualsQualMap) identity.ListUsersRequest {
	request := identity.ListUsersRequest{}

	if equalQuals["external_identifier"] != nil {
		request.ExternalIdentifier = types.String(equalQuals["external_identifier"].GetStringValue())
	}
	if equalQuals["identity_provider_id"] != nil {
		request.IdentityProviderId = types.String(equalQuals["identity_provider_id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = identity.UserLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}
	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}

	return request
}
