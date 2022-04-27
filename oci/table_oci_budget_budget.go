package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/budget"
	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableBudget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_budget_budget",
		Description: "OCI Budget",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBudget,
		},
		List: &plugin.ListConfig{
			Hydrate: listBudgets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The display name of the budget.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the budget.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "amount",
				Description: "The amount of the budget expressed in the currency of the customer's rate card.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "actual_spend",
				Description: "The actual spend in currency for the current budget cycle.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the budget.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_compartment_id",
				Description: "This is DEPRECATED. For backwards compatability, the property will be populated when targetType is COMPARTMENT AND targets contains EXACT ONE target compartment ocid. For all other scenarios, this property will be left empty.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "Time that budget was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "alert_rule_count",
				Description: "Total number of alert rules in the budget.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "budget_processing_period_start_offset",
				Description: "The number of days offset from the first day of the month, at which the budget processing period starts.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "The description of the budget.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "forecasted_spend",
				Description: "The forecasted spend in currency by the end of the current budget cycle.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "reset_period",
				Description: "The reset period for the budget.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_type",
				Description: "The type of target on which the budget is applied.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_spend_computed",
				Description: "The time that the budget spend was last computed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeSpendComputed.Time"),
			},
			{
				Name:        "time_updated",
				Description: "Time that budget was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "version",
				Description: "Version of the budget. Starts from 1 and increments by 1.",
				Type:        proto.ColumnType_INT,
			},

			// json fields
			{
				Name:        "targets",
				Description: "The list of targets on which the budget is applied.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(budgetTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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

func listBudgets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listBudgets", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := budgetService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := budget.ListBudgetsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		TargetType:    "ALL",
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
		request.LifecycleState = budget.ListBudgetsLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BudgetClient.ListBudgets(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, budget := range response.Items {
			d.StreamListItem(ctx, budget)

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

//// HYDRATE FUNCTION

func getBudget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getBudget", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := budgetService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := budget.GetBudgetRequest{
		BudgetId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BudgetClient.GetBudget(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Budget, nil
}

//// TRANSFORM FUNCTION

func budgetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case budget.BudgetSummary:
		budget := d.HydrateItem.(budget.BudgetSummary)
		freeformTags = budget.FreeformTags
		definedTags = budget.DefinedTags
	case budget.Budget:
		budget := d.HydrateItem.(budget.Budget)
		freeformTags = budget.FreeformTags
		definedTags = budget.DefinedTags
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
