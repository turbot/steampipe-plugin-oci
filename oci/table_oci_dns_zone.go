package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/dns"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableDnsZone(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_dns_zone",
		Description: "OCI DNS Zone",
		List: &plugin.ListConfig{
			Hydrate: listDnsZones,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDnsZone,
		},
		GetMatrixItem: BuildCompartmentList,
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
				Name:        "lifecycle_state",
				Description: "The current state of the zone resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone_type",
				Description: "The type of the zone. Must be either `PRIMARY` or `SECONDARY`. `SECONDARY` is only supported for GLOBAL zones.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the zone was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "is_protected",
				Description: "A Boolean flag indicating whether or not parts of the resource are unable to be explicitly managed.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "scope",
				Description: "The scope of the zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "serial",
				Description: "The current serial of the zone. As seen in the zone's SOA record.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "self",
				Description: "The canonical absolute URL of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "Version is the never-repeating, totally-orderable, version of the zone, from which the serial field of the zone's SOA record is derived.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "view_id",
				Description: "The OCID of the private view containing the zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},

			// json fields
			{
				Name:        "external_masters",
				Description: "External master servers for the zone.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDnsZone,
			},
			{
				Name:        "nameservers",
				Description: "The authoritative nameservers for the zone.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDnsZone,
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

func listDnsZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listDnsZones", "Compartment", compartment)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := dnsService(ctx, d)
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
		zones, err := session.DnsClient.ListZones(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, zone := range zones.Items {
			d.StreamListItem(ctx, zone)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				zones.OpcNextPage = nil
			}
		}
		if zones.OpcNextPage != nil {
			request.Page = zones.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDnsZone(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getDnsZone", "Compartment", compartment)

	// Rstrict the api call to only root compartment
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	var id string
	if h.Item != nil {
		id = *h.Item.(dns.ZoneSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	if len(id) == 0 {
		return nil, nil
	}
	// Create Session
	session, err := dnsService(ctx, d)
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

// Priority order for tags
// 2. Defined Tags
// 3. Free-form tags
func dnsZoneTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case dns.Zone:
		zone := d.HydrateItem.(dns.Zone)
		freeformTags = zone.FreeformTags
		definedTags = zone.DefinedTags
	case dns.ZoneSummary:
		zone := d.HydrateItem.(dns.ZoneSummary)
		freeformTags = zone.FreeformTags
		definedTags = zone.DefinedTags
	}

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	if definedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range definedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
