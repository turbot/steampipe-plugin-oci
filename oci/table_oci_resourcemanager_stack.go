package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/resourcemanager"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOciResourceManagerStack(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_resourcemanager_stack",
		Description: "OCI Resource Manager Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getResourceManagerStack,
		},
		List: &plugin.ListConfig{
			Hydrate: listResourceManagerStacks,
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
					Name:    "display_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "Human-readable display name for the stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Unique identifier of the specified stack.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time when the stack was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "General description of the stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stack_drift_status",
				Description: "Drift status of the stack.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getResourceManagerStack,
			},
			{
				Name:        "terraform_version",
				Description: "The version of Terraform specified for the stack.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_drift_last_checked",
				Description: "The date and time when the drift detection was last executed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeDriftLastChecked.Time"),
				Hydrate:     getResourceManagerStack,
			},
			{
				Name:        "config_source",
				Description: "The version of Terraform specified for the stack.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getResourceManagerStack,
			},
			{
				Name:        "variables",
				Description: "Terraform variables associated with this resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getResourceManagerStack,
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

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(resourceManagerStackTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// OCI standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(ociRegionName),
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
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

//// LIST FUNCTION

func listResourceManagerStacks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listResourceManagerStacks", "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := resourceManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := resourcemanager.ListStacksRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}
	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidResourceManagerStackLifecycleState(lifecycleState) {
			request.LifecycleState = resourcemanager.StackLifecycleStateEnum(lifecycleState)
		} else {
			return nil, nil
		}
	}

	// Check for limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ResourceManagerClient.ListStacks(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, resource := range response.Items {
			d.StreamListItem(ctx, resource)

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

//// HYDRATE FUNCTION

func getResourceManagerStack(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getResourceManagerStack", "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(resourcemanager.StackSummary).Id
	} else {
		// Restrict the api call to only root compartment
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := resourceManagerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := resourcemanager.GetStackRequest{
		StackId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ResourceManagerClient.GetStack(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Stack, nil
}

//// TRANSFORM FUNCTION

func resourceManagerStackTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case resourcemanager.StackSummary:
		stacks := d.HydrateItem.(resourcemanager.StackSummary)
		freeformTags = stacks.FreeformTags
		definedTags = stacks.DefinedTags
	case resourcemanager.Stack:
		stack := d.HydrateItem.(resourcemanager.Stack)
		freeformTags = stack.FreeformTags
		definedTags = stack.DefinedTags
	}

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	if definedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range definedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

func isValidResourceManagerStackLifecycleState(state string) bool {
	stateType := resourcemanager.StackLifecycleStateEnum(state)
	switch stateType {
	case resourcemanager.StackLifecycleStateActive, resourcemanager.StackLifecycleStateCreating, resourcemanager.StackLifecycleStateDeleted, resourcemanager.StackLifecycleStateDeleting, resourcemanager.StackLifecycleStateFailed:
		return true
	}
	return false
}
