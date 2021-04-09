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

func tableCoreLocalPeeringGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_local_peering_gateway",
		Description: "OCI Core Local Peering Gateway",
		List: &plugin.ListConfig{
			Hydrate: listPeeringGateway,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getPeeringGateway,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "The LPG's Oracle ID",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN that uses the LPG.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The LPG's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the LPG was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "is_cross_tenancy_peering",
				Description: "Whether the VCN at the other end of the peering is in a different tenancy.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "peer_advertised_cidr",
				Description: "The smallest aggregate CIDR that contains all the CIDR routes advertised by the VCN at the other end of the peering from this LPG.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "peering_status_details",
				Description: "Additional information regarding the peering status.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "route_table_id",
				Description: "The OCID of the route table the LPG is using.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},

			// json fields
			{
				Name:        "peering_status",
				Description: "Whether the LPG is peered with another LPG.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "peer_advertised_cidr_details",
				Description: "The specific ranges of IP addresses available on or via the VCN at the other end of the peering from this LPG.",
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
				Transform:   transform.From(gatewayTags),
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

func listPeeringGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listPeeringGateway", "Compartment", compartment, "VCN", region)

	// Create Session
	session, err := VCNService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListLocalPeeringGatewaysRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		gateways, err := session.VCNClient.ListLocalPeeringGateways(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, gateway := range gateways.Items {
			d.StreamListItem(ctx, gateway)
		}
		if gateways.OpcNextPage != nil {
			request.Page = gateways.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getPeeringGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getPeeringGateway", "Compartment", compartment, "OCI_REGION", region)

	// Rstrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	if len(id) == 0 {
		return nil, nil
	}
	// Create Session
	session, err := VCNService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetLocalPeeringGatewayRequest{
		LocalPeeringGatewayId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VCNClient.GetLocalPeeringGateway(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.LocalPeeringGateway, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 2. Defined Tags
// 3. Free-form tags
func gatewayTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	localPeeringGateway := d.HydrateItem.(core.LocalPeeringGateway)

	var tags map[string]interface{}

	if localPeeringGateway.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range localPeeringGateway.FreeformTags {
			tags[k] = v
		}
	}

	if localPeeringGateway.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range localPeeringGateway.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
