package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVcn(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_vcn",
		Description: "OCI Core VCN",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getVcn,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVcns,
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
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The VCN's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The VCN's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the VCN was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "cidr_block",
				Description: "The first CIDR IP address from cidrBlocks.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "default_dhcp_options_id",
				Description: "The OCID for the VCN's default set of DHCP options.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultDhcpOptionsId"),
			},
			{
				Name:        "default_route_table_id",
				Description: "The OCID of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultRouteTableId"),
			},
			{
				Name:        "default_security_list_id",
				Description: "The OCID for the VCN's default security list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DefaultSecurityListId"),
			},
			{
				Name:        "dns_label",
				Description: "A DNS label for the VCN, used in conjunction with the VNIC's hostname and subnet's DNS label to form a fully qualified domain name (FQDN) for each VNIC within this subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DnsLabel"),
			},
			{
				Name:        "ipv6_cidr_block",
				Description: "For an IPv6-enabled VCN, this is the IPv6 CIDR block for the VCN's private IP address space.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Ipv6CidrBlock"),
			},
			{
				Name:        "ipv6_public_cidr_block",
				Description: "For an IPv6-enabled VCN, this is the IPv6 CIDR block for the VCN's public IP address space.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ipv6PublicCidrBlock"),
			},
			{
				Name:        "vcn_domain_name",
				Description: "The VCN's domain name, which consists of the VCN's DNS label, and the oraclevcn.com domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_blocks",
				Description: "The list of IPv4 CIDR blocks the VCN will use.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ipv6_cidr_blocks",
				Description: "For an IPv6-enabled VCN, this is the list of IPv6 CIDR blocks for the VCN's IP address space. The CIDRs are provided by Oracle and the sizes are always /56.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Ipv6CidrBlocks"),
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
				Transform:   transform.From(vcnTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCoreVcns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreVcns", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListVcnsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Check for additional filters
	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = core.VcnLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListVcns(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, network := range response.Items {
			d.StreamListItem(ctx, network)

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

func getVcn(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVcn")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getVcn", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVcnRequest{
		VcnId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetVcn(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Vcn, nil
}

//// TRANSFORM FUNCTION

func vcnTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	vcn := d.HydrateItem.(core.Vcn)

	var tags map[string]interface{}

	if vcn.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range vcn.FreeformTags {
			tags[k] = v
		}
	}

	if vcn.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range vcn.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
