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

func tableCorePublicIP(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_public_ip",
		Description: "OCI Core Public IP",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getPublicIP,
		},
		List: &plugin.ListConfig{
			Hydrate: listPublicIPs,
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
				Description: "The public IP's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The public IP's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "Whether the public IP is regional or specific to a particular availability domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assigned_entity_id",
				Description: "The OCID of the entity the public IP is assigned to, or in the process of being assigned to.",
				Type:        proto.ColumnType_STRING,
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
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifetime",
				Description: "Defines when the public IP is deleted and released back to Oracle's public IP pool.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip_id",
				Description: "The OCID of the private IP that the public IP is currently assigned to, or in the process of being assigned to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_ip_pool_id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the pool object created in the current tenancy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the publicIP was created.",
				Type:        proto.ColumnType_STRING,
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listPublicIPs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("listPublicIPs", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := virtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListPublicIpsRequest{
		CompartmentId: types.String(compartment),
		Scope:         `REGION`,
		RequestMetadata: oci_common.RequestMetadata{
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

func getPublicIP(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getPublicIP")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("oci.getPublicIP", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := virtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetPublicIpRequest{
		PublicIpId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetPublicIp(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.PublicIp, nil
}

//// TRANSFORM FUNCTION

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
