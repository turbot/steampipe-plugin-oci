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

func tableCoreSecurityList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_security_list",
		Description: "A security list is a virtual firewall for an instance, with ingress and egress rules that specify the types of traffic allowed in and out.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getCoreSecurityList,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreSecurityLists,
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
				Description: "The security list's Oracle Cloud ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the VCN the security list belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The security list's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the security list was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "egress_security_rules",
				Description: "Rules for allowing egress IP packets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ingress_security_rules",
				Description: "Rules for allowing ingress IP packets.",
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
				Transform:   transform.From(securityListTags),
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

func listCoreSecurityLists(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreSecurityLists", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListSecurityListsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.VirtualNetworkClient.ListSecurityLists(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, securityList := range response.Items {
			d.StreamListItem(ctx, securityList)
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

func getCoreSecurityList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreSecurityList")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCoreSecurityList", "Compartment", compartment, "OCI_REGION", region)

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

	request := core.GetSecurityListRequest{
		SecurityListId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetSecurityList(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.SecurityList, nil
}

//// TRANSFORM FUNCTION

func securityListTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	securityList := d.HydrateItem.(core.SecurityList)

	var tags map[string]interface{}

	if securityList.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range securityList.FreeformTags {
			tags[k] = v
		}
	}

	if securityList.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range securityList.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
