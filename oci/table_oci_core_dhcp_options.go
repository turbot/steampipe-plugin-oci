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

func tableCoreDhcpOptions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_dhcp_options",
		Description: "OCI Core DHCP Options",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreDhcpOption,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreDhcpOptions,
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
				Description: "The OCID of the DHCP options.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name of the DHCP options.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The DHCP options's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN the DHCP options belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the DHCP options was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// JSON columns
			{
				Name:        "options",
				Description: "The collection of individual DHCP options.",
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
				Transform:   transform.From(dhcpOptionTags),
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

func listCoreDhcpOptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listCoreDhcpOptions", "Compartment", compartment, "OCI_REGION", region)

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

	request := core.ListDhcpOptionsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListDhcpOptions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, dhcpOption := range response.Items {
			d.StreamListItem(ctx, dhcpOption)

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

func getCoreDhcpOption(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreDhcpOption")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCoreDhcpOption", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetDhcpOptionsRequest{
		DhcpId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetDhcpOptions(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DhcpOptions, nil
}

//// TRANSFORM FUNCTION

func dhcpOptionTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	dhcpOption := d.HydrateItem.(core.DhcpOptions)

	var tags map[string]interface{}

	if dhcpOption.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range dhcpOption.FreeformTags {
			tags[k] = v
		}
	}

	if dhcpOption.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range dhcpOption.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}
		}
	}

	return tags, nil
}
