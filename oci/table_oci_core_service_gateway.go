package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreServiceGateway(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_service_gateway",
		Description: "OCI Core Service Gateway",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreServiceGateway,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreServiceGateways,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "vcn_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the service gateway.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name of the service gateway.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The service gateway's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN the service gateway belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "route_table_id",
				Description: "The OCID of the route table the service gateway is using.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the service gateway was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "block_traffic",
				Description: "Specifies whether the service gateway blocks traffic through it.",
				Type:        proto.ColumnType_BOOL,
			},

			// json columns
			{
				Name:        "services",
				Description: "List of the Service objects enabled for this service gateway.",
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
				Transform:   transform.From(serviceGatewayTags),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCoreServiceGateways(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.listCoreServiceGateways", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListServiceGatewaysRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["vcn_id"] != nil {
		vcnId := equalQuals["vcn_id"].GetStringValue()
		request.VcnId = types.String(vcnId)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = core.ServiceGatewayLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListServiceGateways(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, serviceGateway := range response.Items {
			d.StreamListItem(ctx, serviceGateway)

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

//// HYDRATE FUNCTION

func getCoreServiceGateway(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreServiceGateway")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.getCoreServiceGateway", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetServiceGatewayRequest{
		ServiceGatewayId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.VirtualNetworkClient.GetServiceGateway(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.ServiceGateway, nil
}

//// TRANSFORM FUNCTION

func serviceGatewayTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	serviceGateway := d.HydrateItem.(core.ServiceGateway)

	var tags map[string]interface{}

	if serviceGateway.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range serviceGateway.FreeformTags {
			tags[k] = v
		}
	}

	if serviceGateway.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range serviceGateway.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}
		}
	}

	return tags, nil
}
