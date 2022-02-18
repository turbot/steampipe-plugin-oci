package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_group",
		Description: "OCI Identity Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listGroup,
			KeyColumns: []*plugin.KeyColumn{
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
				Description: "The name assign to the group during creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "The OCID of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The group's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inactive_status",
				Description: "The detailed status of INACTIVE lifecycleState.",
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
				Transform:   transform.From(groupTags),
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

func listGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	equalQuals := d.KeyColumnQuals

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListGroupsRequest{
		CompartmentId: &session.TenancyID,
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Check for additional filters
	if equalQuals["name"] != nil {
		name := equalQuals["name"].GetStringValue()
		request.Name = types.String(name)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = identity.GroupLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListGroups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, group := range response.Items {
			d.StreamListItem(ctx, group)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func getGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getGroup")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetGroupRequest{
		GroupId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetGroup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Group, nil
}

//// TRANSFORM FUNCTION

func groupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	group := d.HydrateItem.(identity.Group)

	var tags map[string]interface{}

	if group.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range group.FreeformTags {
			tags[k] = v
		}
	}

	if group.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range group.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
