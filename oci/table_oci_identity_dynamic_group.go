package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityDynamicGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_dynamic_group",
		Description: "OCI Identity Dynamic Group",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIdentityDynamicGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listIdentityDynamicGroups,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
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
				Description: "The name you assign to the group during creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "id",
				Description: "The OCID of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The group's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The description you assign to the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inactive_status",
				Description: "The detailed status of INACTIVE lifecycleState.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "matching_rule",
				Description: "A rule string that defines which instance certificates will be matched.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the group was created, in the format defined by RFC3339.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Transform:   transform.From(dynamicGroupTags),
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

func listIdentityDynamicGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listIdentityDynamicGroups", "Compartment", compartment)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListDynamicGroupsRequest{
		CompartmentId: types.String(compartment),
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
		request.LifecycleState = identity.DynamicGroupLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListDynamicGroups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, dynamicGroup := range response.Items {
			d.StreamListItem(ctx, dynamicGroup)

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

func getIdentityDynamicGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getIdentityDynamicGroup")

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty dynamic group id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetDynamicGroupRequest{
		DynamicGroupId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetDynamicGroup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DynamicGroup, nil
}

//// TRANSFORM FUNCTION

func dynamicGroupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dynanicGroup := d.HydrateItem.(identity.DynamicGroup)

	var tags map[string]interface{}

	if dynanicGroup.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range dynanicGroup.FreeformTags {
			tags[k] = v
		}
	}

	if dynanicGroup.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range dynanicGroup.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
