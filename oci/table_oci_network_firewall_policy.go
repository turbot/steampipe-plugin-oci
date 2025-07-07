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
				Name:        "time_created",
				Description: "Time that Network Firewall Policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "is_firewall_attached",
				Description: "To determine if any Network Firewall is associated with this Network Firewall Policy.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getNetworkFirewallPolicy,
				Transform:   transform.From(isFirewallAttached),
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
				Name:        "application_lists",
				Description: "A mapping of strings to arrays of Application objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyApplications,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "decryption_profiles",
				Description: "A mapping of strings to DecryptionProfile objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyDecryptionProfiles,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "decryption_rules",
				Description: "List of Decryption Rules defining the behavior of the policy. The first rule with a matching condition determines the action taken upon network traffic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyDecryptionRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "ip_address_lists",
				Description: "Map defining IP address lists of the policy. The value of an entry is a list of IP addresses or prefixes in CIDR notation. The associated key is the identifier by which the IP address list is referenced.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyAddresses,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "mapped_secrets",
				Description: "A mapping of strings to MappedSecret objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyMappedSecrets,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "security_rules",
				Description: "List of Security Rules defining the behavior of the policy. The first rule with a matching condition determines the action taken upon network traffic.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicySecurityRules,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "url_lists",
				Description: "A mapping of strings to arrays of UrlPattern objects.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listNetworkFirewallPolicyURLs,
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listNetworkFirewallPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicies", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := networkFirewallService(ctx, d, region)
	if err != nil {
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicies", "connection_error", err)
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
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicies", "api_error", err)
			return nil, err
		}
		for _, firewallPolicy := range response.Items {
			d.StreamListItem(ctx, firewallPolicy)

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

func getNetworkFirewallPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.getNetworkFirewallPolicy", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
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
		logger.Error("oci_network_firewall_policy.getNetworkFirewallPolicy", "connection_error", err)
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
		logger.Error("oci_network_firewall_policy.getNetworkFirewallPolicy", "api_error", err)
		return nil, err
	}
	return response.NetworkFirewallPolicy, nil
}

func listNetworkFirewallPolicyApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyApplications", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyApplications", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListApplicationsRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var applications []networkfirewall.ApplicationSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListApplications(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyApplications", "api_error", err)
			return nil, err
		}
		applications = append(applications, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return applications, nil
}

func listNetworkFirewallPolicyDecryptionProfiles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionProfiles", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionProfiles", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListDecryptionProfilesRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var decryptionProfiles []networkfirewall.DecryptionProfileSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListDecryptionProfiles(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionProfiles", "api_error", err)
			return nil, err
		}
		decryptionProfiles = append(decryptionProfiles, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return decryptionProfiles, nil
}

func listNetworkFirewallPolicyDecryptionRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionRules", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionRules", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListDecryptionRulesRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}


	var decryptionRules []networkfirewall.DecryptionRuleSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListDecryptionRules(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyDecryptionRules", "api_error", err)
			return nil, err
		}
		decryptionRules = append(decryptionRules, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return decryptionRules, nil
}

func listNetworkFirewallPolicyMappedSecrets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyMappedSecrets", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyMappedSecrets", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListMappedSecretsRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var mappedSecrets []networkfirewall.MappedSecretSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListMappedSecrets(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyMappedSecrets", "api_error", err)
			return nil, err
		}
		mappedSecrets = append(mappedSecrets, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return mappedSecrets, nil
}

func listNetworkFirewallPolicySecurityRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicySecurityRules", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicySecurityRules", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListSecurityRulesRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var securityRules []networkfirewall.SecurityRuleSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListSecurityRules(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicySecurityRules", "api_error", err)
			return nil, err
		}
		securityRules = append(securityRules, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return securityRules, nil
}

func listNetworkFirewallPolicyAddresses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyAddresses", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyAddresses", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListAddressListsRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var addressLists []networkfirewall.AddressListSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListAddressLists(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyAddresses", "api_error", err)
			return nil, err
		}
		addressLists = append(addressLists, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return addressLists, nil
}

func listNetworkFirewallPolicyURLs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_network_firewall_policy.listNetworkFirewallPolicyURLs", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		switch h.Item.(type) {
		case networkfirewall.NetworkFirewallPolicySummary:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicySummary).Id
		case networkfirewall.NetworkFirewallPolicy:
			id = *h.Item.(networkfirewall.NetworkFirewallPolicy).Id
		}
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
		logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyURLs", "connection_error", err)
		return nil, err
	}

	request := networkfirewall.ListUrlListsRequest{
		NetworkFirewallPolicyId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	var urlLists []networkfirewall.UrlListSummary
	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkFirewallClient.ListUrlLists(ctx, request)
		if err != nil {
			logger.Error("oci_network_firewall_policy.listNetworkFirewallPolicyURLs", "api_error", err)
			return nil, err
		}
		urlLists = append(urlLists, response.Items...)
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return urlLists, nil
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

func isFirewallAttached(_ context.Context, d *transform.TransformData) (interface{}, error) {
	firewallPolicy := d.HydrateItem.(networkfirewall.NetworkFirewallPolicy)

	if firewallPolicy.AttachedNetworkFirewallCount != nil && *firewallPolicy.AttachedNetworkFirewallCount > 0 {
		return true, nil
	}

	return false, nil
}

// Build additional filters
func buildNetworkFirewallPolicyFilters(equalQuals plugin.KeyColumnEqualsQualMap) networkfirewall.ListNetworkFirewallPoliciesRequest {
	request := networkfirewall.ListNetworkFirewallPoliciesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = networkfirewall.ListNetworkFirewallPoliciesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
