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

func tableWafWebAppFirewall(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_waf_web_app_firewall",
		Description:      "OCI WAF Web Application Firewall",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getWafWebAppFirewall,
		},
		List: &plugin.ListConfig{
			Hydrate: listWafWebAppFirewalls,
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
					Name:    "web_app_firewall_policy_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the WebAppFirewall.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "display_name",
				Description: "WebAppFirewall display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "backend_type",
				Description: "Type of the WebAppFirewall, e.g. LOAD_BALANCER.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafWebAppFirewall,
				Transform:   transform.FromField("BackendType"),
			},
			{
				Name:        "web_app_firewall_policy_id",
				Description: "The OCID of the WebAppFirewallPolicy attached to this firewall.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WebAppFirewallPolicyId"),
			},
			{
				Name:        "load_balancer_id",
				Description: "The OCID of the LoadBalancer to which the policy is attached (LOAD_BALANCER backend type).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafWebAppFirewall,
				Transform:   transform.FromField("LoadBalancerId"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the WebAppFirewall.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LifecycleState"),
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LifecycleDetails"),
			},
			{
				Name:        "time_created",
				Description: "The time the WebAppFirewall was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the WebAppFirewall was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefinedTags"),
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FreeformTags"),
			},
			{
				Name:        "system_tags",
				Description: "Usage of system tag keys.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SystemTags"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(wafWebAppFirewallTags),
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

func listWafWebAppFirewalls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_web_app_firewall.listWafWebAppFirewalls", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	session, err := wafService(ctx, d, region)
	if err != nil {
		logger.Error("oci_waf_web_app_firewall.listWafWebAppFirewalls", "connection_error", err)
		return nil, err
	}

	request := waf.ListWebAppFirewallsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["web_app_firewall_policy_id"] != nil {
		request.WebAppFirewallPolicyId = types.String(equalQuals["web_app_firewall_policy_id"].GetStringValue())
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.WafClient.ListWebAppFirewalls(ctx, request)
		if err != nil {
			logger.Error("oci_waf_web_app_firewall.listWafWebAppFirewalls", "api_error", err)
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

func getWafWebAppFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_web_app_firewall.getWafWebAppFirewall", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		summary := h.Item.(waf.WebAppFirewallSummary)
		id = *summary.GetId()
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
		logger.Error("oci_waf_web_app_firewall.getWafWebAppFirewall", "connection_error", err)
		return nil, err
	}

	request := waf.GetWebAppFirewallRequest{
		WebAppFirewallId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.WafClient.GetWebAppFirewall(ctx, request)
	if err != nil {
		logger.Error("oci_waf_web_app_firewall.getWafWebAppFirewall", "api_error", err)
		return nil, err
	}

	return response.WebAppFirewall, nil
}

//// TRANSFORM FUNCTIONS

func wafWebAppFirewallTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch item := d.HydrateItem.(type) {
	case waf.WebAppFirewallLoadBalancerSummary:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	case waf.WebAppFirewallLoadBalancer:
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
