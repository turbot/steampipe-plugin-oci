package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/budget"
	"github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableBudgetAlertRule(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_budget_alert_rule",
		Description: "OCI Budget Alert Rule",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id,budget_id"}),
			Hydrate:    getBudgetAlertRule,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listBudgets,
			Hydrate:       listBudgetAlertRules,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The name of the alert rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the alert rule.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "budget_id",
				Description: "The OCID of the budget",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "threshold",
				Description: "The threshold for triggering the alert. If thresholdType is PERCENTAGE, the maximum value is 10000.",
				Type:        proto.ColumnType_DOUBLE,
			},
			{
				Name:        "threshold_type",
				Description: "The type of threshold.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the alert rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that budget was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "description",
				Description: "The description of the alert rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "message",
				Description: "Custom message sent when alert is triggered.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recipients",
				Description: "Delimited list of email addresses to receive the alert when it triggers.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "Time that budget was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "Type",
				Description: "The type of alert.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "Version of the alert rule. Starts from 1 and increments by 1.",
				Type:        proto.ColumnType_INT,
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
				Transform:   transform.From(alertRuleTags),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listBudgetAlertRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listBudgets", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := budgetService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	id := h.Item.(budget.BudgetSummary).Id

	request := budget.ListAlertRulesRequest{
		BudgetId: id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BudgetClient.ListAlertRules(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, rule := range response.Items {
			d.StreamListItem(ctx, rule)
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

func getBudgetAlertRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getBudgetAlertRule", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	budgetId := d.KeyColumnQuals["budget_id"].GetStringValue()
	ruleId := d.KeyColumnQuals["rule_id"].GetStringValue()

	// handle empty id in get call
	if budgetId == "" || ruleId == "" {
		return nil, nil
	}

	// Create Session
	session, err := budgetService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := budget.GetAlertRuleRequest{
		BudgetId:    types.String(budgetId),
		AlertRuleId: types.String(ruleId),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.BudgetClient.GetAlertRule(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.AlertRule, nil
}

//// TRANSFORM FUNCTION

func alertRuleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case budget.BudgetSummary:
		alertRule := d.HydrateItem.(budget.AlertRuleSummary)
		freeformTags = alertRule.FreeformTags
		definedTags = alertRule.DefinedTags
	case budget.AlertRule:
		alertRule := d.HydrateItem.(budget.AlertRule)
		freeformTags = alertRule.FreeformTags
		definedTags = alertRule.DefinedTags
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
