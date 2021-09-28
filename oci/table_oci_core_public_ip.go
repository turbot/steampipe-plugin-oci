package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCorePublicIP(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_public_ip",
		Description: "OCI Core Public IP",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCorePublicIP,
		},
		List: &plugin.ListConfig{
			Hydrate: listCorePublicIPs,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The public IP's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "assigned_entity_id",
				Description: "The OCID of the entity the public IP is assigned to, or in the process of being assigned to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The public IP's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the public IP was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "assigned_entity_type",
				Description: "The type of entity the public IP is assigned to, or in the process of being assigned to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The public IP's availability domain. This property is set only for ephemeral public IPs that are assigned to a private IP.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The public IP address of the publicIp object.",
				Type:        proto.ColumnType_IPADDR,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifetime",
				Description: "Defines when the public IP is deleted and released back to Oracle's public IP pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_pool_id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the pool object created in the current tenancy.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "scope",
				Description: "Indicates whether the public IP is regional or specific to a particular availability domain.",
				Type:        proto.ColumnType_STRING,
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
				Transform:   transform.From(publicIPTags),
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

func listCorePublicIPs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	plugin.Logger(ctx).Error("listCorePublicIPs", "Compartment", compartment, "OCI_REGION", region)

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

	request := core.ListPublicIpsRequest{
		CompartmentId: types.String(compartment),
		Scope:         "REGION",
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListPublicIps(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, ip := range response.Items {
			d.StreamListItem(ctx, ip)

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

//// HYDRATE FUNCTIONS

func getCorePublicIP(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	plugin.Logger(ctx).Error("getCorePublicIP", "Compartment", compartment, "OCI_REGION", region)

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

	request := core.GetPublicIpRequest{
		PublicIpId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetPublicIp(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.PublicIp, nil
}

//// TRANSFORM FUNCTIONS

func publicIPTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	publicIP := d.HydrateItem.(core.PublicIp)

	var tags map[string]interface{}

	if publicIP.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range publicIP.FreeformTags {
			tags[k] = v
		}
	}

	if publicIP.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range publicIP.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}
		}
	}

	return tags, nil
}
