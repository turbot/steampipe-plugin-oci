package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreSubnet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_subnet",
		Description: "OCI Core Subnet",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getCoreSubnet,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreSubnets,
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
				Description: "The subnet's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN the subnet is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The subnet's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "route_table_id",
				Description: "The OCID of the route table that the subnet uses.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the subnet was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "availability_domain",
				Description: "The subnet's availability domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cidr_block",
				Description: "The subnet's CIDR block.",
				Type:        proto.ColumnType_CIDR,
			},
			{
				Name:        "dhcp_options_id",
				Description: "The OCID of the set of DHCP options that the subnet uses.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "dns_label",
				Description: "A DNS label for the subnet, used in conjunction with the VNIC's hostname and VCN's DNS label to form a fully qualified domain name (FQDN) for each VNIC within this subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "ipv6_cidr_block",
				Description: "For an IPv6-enabled subnet, this is the IPv6 CIDR block for the subnet's private IP address space.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Ipv6CidrBlock"),
			},
			{
				Name:        "ipv6_public_cidr_block",
				Description: "For an IPv6-enabled subnet, this is the IPv6 CIDR block for the subnet's public IP address space.",
				Type:        proto.ColumnType_CIDR,
				Transform:   transform.FromField("Ipv6PublicCidrBlock"),
			},
			{
				Name:        "ipv6_virtual_router_ip",
				Description: "For an IPv6-enabled subnet, this is the IPv6 address of the virtual router.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromField("Ipv6VirtualRouterIp"),
			},
			{
				Name:        "prohibit_public_ip_on_vnic",
				Description: "Indicates whether VNICs within this subnet can have public IP addresses.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "subnet_domain_name",
				Description: "The subnet's domain name, which consists of the subnet's DNS label, the VCN's DNS label, and the `oraclevcn.com` domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "virtual_router_ip",
				Description: "The IP address of the virtual router.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "virtual_router_mac",
				Description: "The MAC address of the virtual router.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "security_list_ids",
				Description: "The OCIDs of the security list or lists that the subnet uses.",
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
				Transform:   transform.From(subnetTags),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCoreSubnets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreSubnets", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListSubnetsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListSubnets(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, subnet := range response.Items {
			d.StreamListItem(ctx, subnet)
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

func getCoreSubnet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreSubnet")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCoreSubnet", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetSubnetRequest{
		SubnetId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetSubnet(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Subnet, nil
}

//// TRANSFORM FUNCTION

func subnetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	subnet := d.HydrateItem.(core.Subnet)

	var tags map[string]interface{}

	if subnet.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range subnet.FreeformTags {
			tags[k] = v
		}
	}

	if subnet.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range subnet.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
