package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreInternetGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_internet_gateway",
		Description: "OCI Core Internet Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getCoreInternetGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreInternetGateways,
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
				Name:        "id",
				Description: "The internet gateway's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN the internet gateway belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "is_enabled",
				Description: "Whether the gateway is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_state",
				Description: "The internet gateway's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the internet gateway was created.",
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
				Transform:   transform.From(internetGatewayTags),
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

type internetGatewayInfo struct {
	core.InternetGateway
	Region string
}

//// LIST FUNCTION

func listCoreInternetGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreInternetGateways", "Compartment", compartment, "OCI_REGION", region)

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

	request := core.ListInternetGatewaysRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListInternetGateways(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, internetGateway := range response.Items {
			d.StreamListItem(ctx, internetGatewayInfo{internetGateway, region})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				response.OpcNextPage = nil
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

func getCoreInternetGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreInternetGateway")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreInternetGateway", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty internet gateway id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetInternetGatewayRequest{
		IgId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetInternetGateway(ctx, request)
	if err != nil {
		return nil, err
	}

	return internetGatewayInfo{response.InternetGateway, region}, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func internetGatewayTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	internetGateway := d.HydrateItem.(internetGatewayInfo).InternetGateway

	var tags map[string]interface{}

	if internetGateway.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range internetGateway.FreeformTags {
			tags[k] = v
		}
	}

	if internetGateway.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range internetGateway.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
