package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/networkfirewall"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strings"
)

//// TABLE DEFINITION

func tableNetworkFirewallPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_network_firewall_policy",
		Description:      "OCI Network Firewall Policy",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getNetworkFirewallPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewallPolicies,
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
					Name:    "id",
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
				Description: "The OCID of the Network Firewall Policy resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name for the Network Firewall Policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_lists",
				Description: "A mapping of strings to arrays of Application objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "decryption_profiles",
				Description: "A mapping of strings to DecryptionProfile objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "decryption_rules",
				Description: "List of Decryption Rules defining the behavior of the policy. The first rule with a matching condition determines the action taken upon network traffic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "ip_address_lists",
				Description: "Map defining IP address lists of the policy. The value of an entry is a list of IP addresses or prefixes in CIDR notation. The associated key is the identifier by which the IP address list is referenced.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "is_firewall_attached",
				Description: "To determine if any Network Firewall is associated with this Network Firewall Policy.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Network Firewall.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mapped_secrets",
				Description: "A mapping of strings to MappedSecret objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "security_rules",
				Description: "List of Security Rules defining the behavior of the policy. The first rule with a matching condition determines the action taken upon network traffic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "url_lists",
				Description: "A mapping of strings to arrays of UrlPattern objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNetworkFirewallPolicy,
			},
			{
				Name:        "time_created",
				Description: "Time that Network Firewall Policy was created.",
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
				Transform:   transform.From(networkFirewallPolicyTags),
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

func listNetworkFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listNetworkFirewallPolicies", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := networkFirewallService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildNetworkFirewallPolicyFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(100)
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
		response, err := session.NetworkFirewallClient.ListNetworkFirewallPolicies(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, firewallPolicy := range response.Items {
			d.StreamListItem(ctx, firewallPolicy)

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

func getNetworkFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getNetworkFirewallPolicy", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := networkFirewallService(ctx, d, region)
	if err != nil {
		logger.Error("getNetworkFirewallPolicy", "error_NetworkFirewallService", err)
		return nil, err
	}

	request := networkfirewall.GetNetworkFirewallPolicyRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.NetworkFirewallClient.GetNetworkFirewallPolicy(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.NetworkFirewallPolicy, nil
}

//// TRANSFORM FUNCTION

func networkFirewallPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case networkfirewall.NetworkFirewallPolicy:
		firewallPolicy := d.HydrateItem.(networkfirewall.NetworkFirewallPolicy)
		freeformTags = firewallPolicy.FreeformTags
		definedTags = firewallPolicy.DefinedTags
	case networkfirewall.NetworkFirewallPolicySummary:
		firewallPolicy := d.HydrateItem.(networkfirewall.NetworkFirewallPolicySummary)
		freeformTags = firewallPolicy.FreeformTags
		definedTags = firewallPolicy.DefinedTags
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

// Build additional filters
func buildNetworkFirewallPolicyFilters(equalQuals plugin.KeyColumnEqualsQualMap) networkfirewall.ListNetworkFirewallPoliciesRequest {
	request := networkfirewall.ListNetworkFirewallPoliciesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["id"] != nil {
		request.Id = types.String(equalQuals["id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = networkfirewall.ListNetworkFirewallPoliciesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
