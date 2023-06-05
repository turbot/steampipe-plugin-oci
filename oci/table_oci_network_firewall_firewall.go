package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/networkfirewall"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableNetworkFirewall(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_network_firewall_firewall",
		Description:      "OCI Network Firewall",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getNetworkFirewall,
		},
		List: &plugin.ListConfig{
			Hydrate: listNetworkFirewalls,
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
					Name:    "network_firewall_policy_id",
					Require: plugin.Optional,
				},
				{
					Name:    "availability_domain",
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
				Description: "The OCID of the Network Firewall resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name for the Network Firewall.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "A filter to return only resources that are present within the specified availability domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Network Firewall was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet associated with the Network Firewall.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipv4_address",
				Description: "IPv4 address for the Network Firewall.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "ipv6_address",
				Description: "IPv6 address for the Network Firewall.",
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Network Firewall.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_firewall_policy_id",
				Description: "The OCID of the Network Firewall Policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "network_security_group_ids",
				Description: "An array of network security groups OCID associated with the Network Firewall.",
				Hydrate:     getNetworkFirewall,
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
				Transform:   transform.From(networkFirewallTags),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listNetworkFirewalls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_firewall.listNetworkFirewalls", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := networkFirewallService(ctx, d, region)
	if err != nil {
		logger.Error("oci_network_firewall_firewall.listNetworkFirewalls", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildNetworkFirewallFilters(equalQuals)
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
		response, err := session.NetworkFirewallClient.ListNetworkFirewalls(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_firewall.listNetworkFirewalls", "api_error", err)
			return nil, err
		}
		for _, firewall := range response.Items {
			d.StreamListItem(ctx, firewall)

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

func getNetworkFirewall(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_firewall.getNetworkFirewall", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(networkfirewall.NetworkFirewallSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
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
		logger.Error("oci_network_firewall_firewall.getNetworkFirewall", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.GetNetworkFirewallRequest{
		NetworkFirewallId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.NetworkFirewallClient.GetNetworkFirewall(ctx, request)
	if err != nil {
		logger.Error("oci_network_firewall_firewall.getNetworkFirewall", "api_error", err)
		return nil, err
	}
	return response.NetworkFirewall, nil
}

//// TRANSFORM FUNCTIONS

func networkFirewallTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case networkfirewall.NetworkFirewall:
		firewall := d.HydrateItem.(networkfirewall.NetworkFirewall)
		freeformTags = firewall.FreeformTags
		definedTags = firewall.DefinedTags
	case networkfirewall.NetworkFirewallSummary:
		firewall := d.HydrateItem.(networkfirewall.NetworkFirewallSummary)
		freeformTags = firewall.FreeformTags
		definedTags = firewall.DefinedTags
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
func buildNetworkFirewallFilters(equalQuals plugin.KeyColumnEqualsQualMap) networkfirewall.ListNetworkFirewallsRequest {
	request := networkfirewall.ListNetworkFirewallsRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["network_firewall_policy_id"] != nil {
		request.NetworkFirewallPolicyId = types.String(equalQuals["network_firewall_policy_id"].GetStringValue())
	}
	if equalQuals["availability_domain"] != nil {
		request.AvailabilityDomain = types.String(equalQuals["availability_domain"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = networkfirewall.ListNetworkFirewallsLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
