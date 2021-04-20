package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/dns"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableDnsZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_dns_zone",
		Description: "OCI Dns Zone",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    getDnsZone,
		},
		List: &plugin.ListConfig{
			Hydrate: listDnsZones,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "zone_type",
				Description: "The type of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_tate",
				Description: "The current state of the zone resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the zone was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "self",
				Description: "The canonical absolute URL of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "The scope of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "Version of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "serial",
				Description: "The current serial of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_protected",
				Description: "A Boolean flag indicating whether or not parts of the resource are unable to be explicitly managed.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DefaultSecurityListId"),
			},
			{
				Name:        "view_id",
				Description: "The OCID of the private view containing the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DnsLabel"),
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
				Transform:   transform.From(dnsZoneTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

func listDnsZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("dns.listDnsZones", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := dnsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := dns.ListZonesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DnsClient.ListZones(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, zone := range response.Items {
			d.StreamListItem(ctx, zone)
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

func getDnsZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVcn")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("dns.getDnsZone", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := dnsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := dns.GetZoneRequest{
		ZoneNameOrId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.DnsClient.GetZone(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Zone, nil
}

//// TRANSFORM FUNCTION

func dnsZoneTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	zone := d.HydrateItem.(dns.ZoneSummary)

	var tags map[string]interface{}

	if zone.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range zone.FreeformTags {
			tags[k] = v
		}
	}

	if zone.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range zone.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
