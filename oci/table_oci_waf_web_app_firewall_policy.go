package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/waf"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableWafWebAppFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_waf_web_app_firewall_policy",
		Description:      "OCI WAF Web Application Firewall Policy",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getWafWebAppFirewallPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listWafWebAppFirewallPolicies,
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
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the WebAppFirewallPolicy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "WebAppFirewallPolicy display name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the WebAppFirewallPolicy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the WebAppFirewallPolicy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the WebAppFirewallPolicy was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "actions",
				Description: "Predefined actions for use in multiple different rules.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
			},
			{
				Name:        "request_access_control",
				Description: "Module that inspects incoming HTTP request access control.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
			},
			{
				Name:        "request_rate_limiting",
				Description: "Module that inspects incoming HTTP request rate limiting.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
			},
			{
				Name:        "request_protection",
				Description: "Module that inspects incoming HTTP request protection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
			},
			{
				Name:        "response_access_control",
				Description: "Module that inspects outgoing HTTP response access control.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
			},
			{
				Name:        "response_protection",
				Description: "Module that inspects outgoing HTTP response protection.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafWebAppFirewallPolicy,
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
			{
				Name:        "system_tags",
				Description: "Usage of system tag keys.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(wafWebAppFirewallPolicyTags),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listWafWebAppFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_web_app_firewall_policy.listWafWebAppFirewallPolicies", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	session, err := wafService(ctx, d, region)
	if err != nil {
		logger.Error("oci_waf_web_app_firewall_policy.listWafWebAppFirewallPolicies", "connection_error", err)
		return nil, err
	}

	request := waf.ListWebAppFirewallPoliciesRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = []waf.WebAppFirewallPolicyLifecycleStateEnum{
			waf.WebAppFirewallPolicyLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue()),
		}
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.WafClient.ListWebAppFirewallPolicies(ctx, request)
		if err != nil {
			logger.Error("oci_waf_web_app_firewall_policy.listWafWebAppFirewallPolicies", "api_error", err)
			return nil, err
		}
		for _, item := range response.Items {
			d.StreamListItem(ctx, item)
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

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getWafWebAppFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_web_app_firewall_policy.getWafWebAppFirewallPolicy", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(waf.WebAppFirewallPolicySummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	if id == "" {
		return nil, nil
	}

	session, err := wafService(ctx, d, region)
	if err != nil {
		logger.Error("oci_waf_web_app_firewall_policy.getWafWebAppFirewallPolicy", "connection_error", err)
		return nil, err
	}

	request := waf.GetWebAppFirewallPolicyRequest{
		WebAppFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.WafClient.GetWebAppFirewallPolicy(ctx, request)
	if err != nil {
		logger.Error("oci_waf_web_app_firewall_policy.getWafWebAppFirewallPolicy", "api_error", err)
		return nil, err
	}

	return response.WebAppFirewallPolicy, nil
}

//// TRANSFORM FUNCTIONS

func wafWebAppFirewallPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch item := d.HydrateItem.(type) {
	case waf.WebAppFirewallPolicySummary:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	case waf.WebAppFirewallPolicy:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
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
