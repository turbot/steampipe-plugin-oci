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

func tableCoreVirtuaCircuit(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_virtual_circuit",
		Description: "OCI Core Virtual Circuit",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getVirtualCircuit,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVirtualCircuits,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The virtual circuit's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},

			{
				Name:        "display_name",
				Description: "A user-friendly name of virtual circuit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the virtual circuit was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Specifies whether the virtual circuit supports private or public peering.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bandwidth_shape_name",
				Description: "The provisioned data rate of the connection. To get a list of the available bandwidth levels (that is, shapes).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of a volume.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "customer_asn",
				Description: "The BGP ASN of the network at the other end of the BGP session from Oracle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gateway_id",
				Description: "The date and time the instance was created",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "oracle_bgp_asn",
				Description: "The date and time the instance was created",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_name",
				Description: "The date and time the instance was created",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_service_id",
				Description: "The OCID of the service offered by the provider (if the customer is connecting via a provider).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProviderServiceId"),
			},
			{
				Name:        "provider_service_key_name",
				Description: "The service key name offered by the provider (if the customer is connecting via a provider).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_service_name",
				Description: "The date and time the instance was created",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reference_comment",
				Description: "Provider-supplied reference information about this virtual circuit(if the customer is connecting via a provider).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_management",
				Description: "Options for the Oracle Cloud Agent software running on the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bgp_session_state",
				Description: "The state of the BGP session associated with the virtual circuit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provider_state",
				Description: "The provider's state in relation to this virtual circuit (if the customer is connecting via a provider).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_type",
				Description: "Provider service type.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "cross_connect_mappings",
				Description: "An array of mappings, each containing properties for a cross-connect or cross-connect group that is associated with this virtual circuit.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "public_prefixes",
				Description: "For a public virtual circuit. The public IP prefixes (CIDRs) the customer wants to advertise across the connection. All prefix sizes are allowed.",
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
				Transform:   transform.From(virtualCircuitTags),
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
				Transform:   transform.FromCamel().Transform(regionName),
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

func listCoreVirtualCircuits(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("core.listCoreVirtualCircuits", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListVirtualCircuitsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListVirtualCircuits(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, circuits := range response.Items {
			d.StreamListItem(ctx, circuits)
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

func getVirtualCircuit(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVirtualCircuit")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("oci.getVirtualCircuit", "Compartment", compartment, "OCI_REGION", region)

	// Rstrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVirtualCircuitRequest{
		VirtualCircuitId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetVirtualCircuit(ctx, request)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Trace("Response--->>>", response.VirtualCircuit)

	return response.VirtualCircuit, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func virtualCircuitTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	circuit := d.HydrateItem.(core.VirtualCircuit)

	var tags map[string]interface{}

	if circuit.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range circuit.FreeformTags {
			tags[k] = v
		}
	}

	if circuit.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range circuit.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
