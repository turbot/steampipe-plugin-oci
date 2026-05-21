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

func tableWafNetworkAddressList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_waf_network_address_list",
		Description:      "OCI WAF Network Address List",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getWafNetworkAddressList,
		},
		List: &plugin.ListConfig{
			Hydrate: listWafNetworkAddressLists,
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
				Description: "The OCID of the NetworkAddressList.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "display_name",
				Description: "NetworkAddressList display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "type",
				Description: "Type of NetworkAddressList (ADDRESSES or VCN_ADDRESSES).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getWafNetworkAddressList,
				Transform:   transform.From(wafNetworkAddressListType),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the NetworkAddressList.",
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
				Description: "The time the NetworkAddressList was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the NetworkAddressList was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "addresses",
				Description: "List of IP address prefixes in CIDR notation (ADDRESSES type).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafNetworkAddressList,
				Transform:   transform.From(wafNetworkAddressListAddresses),
			},
			{
				Name:        "vcn_addresses",
				Description: "List of private address ranges on the OCI VCN (VCN_ADDRESSES type).",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getWafNetworkAddressList,
				Transform:   transform.From(wafNetworkAddressListVcnAddresses),
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
				Transform:   transform.From(wafNetworkAddressListTags),
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

func listWafNetworkAddressLists(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_network_address_list.listWafNetworkAddressLists", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	session, err := wafService(ctx, d, region)
	if err != nil {
		logger.Error("oci_waf_network_address_list.listWafNetworkAddressLists", "connection_error", err)
		return nil, err
	}

	request := waf.ListNetworkAddressListsRequest{
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
		request.LifecycleState = []waf.NetworkAddressListLifecycleStateEnum{
			waf.NetworkAddressListLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue()),
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
		response, err := session.WafClient.ListNetworkAddressLists(ctx, request)
		if err != nil {
			logger.Error("oci_waf_network_address_list.listWafNetworkAddressLists", "api_error", err)
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

func getWafNetworkAddressList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_waf_network_address_list.getWafNetworkAddressList", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		summary := h.Item.(waf.NetworkAddressListSummary)
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
		logger.Error("oci_waf_network_address_list.getWafNetworkAddressList", "connection_error", err)
		return nil, err
	}

	request := waf.GetNetworkAddressListRequest{
		NetworkAddressListId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.WafClient.GetNetworkAddressList(ctx, request)
	if err != nil {
		logger.Error("oci_waf_network_address_list.getWafNetworkAddressList", "api_error", err)
		return nil, err
	}

	return response.NetworkAddressList, nil
}

//// TRANSFORM FUNCTIONS

func wafNetworkAddressListTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch item := d.HydrateItem.(type) {
	case waf.NetworkAddressListAddressesSummary:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	case waf.NetworkAddressListVcnAddressesSummary:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	case waf.NetworkAddressListAddresses:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	case waf.NetworkAddressListVcnAddresses:
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

func wafNetworkAddressListType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	switch d.HydrateItem.(type) {
	case waf.NetworkAddressListAddresses:
		return "ADDRESSES", nil
	case waf.NetworkAddressListVcnAddresses:
		return "VCN_ADDRESSES", nil
	}
	return nil, nil
}

func wafNetworkAddressListAddresses(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if item, ok := d.HydrateItem.(waf.NetworkAddressListAddresses); ok {
		return item.Addresses, nil
	}
	return nil, nil
}

func wafNetworkAddressListVcnAddresses(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if item, ok := d.HydrateItem.(waf.NetworkAddressListVcnAddresses); ok {
		return item.VcnAddresses, nil
	}
	return nil, nil
}
