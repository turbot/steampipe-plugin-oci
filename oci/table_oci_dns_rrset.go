package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/dns"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableDnsRecordSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_dns_rrset",
		Description: "OCI DNS Record Set",
		List: &plugin.ListConfig{
			ParentHydrate: listDnsZones,
			Hydrate:       listDnsRecordSets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "domain",
				Description: "The fully qualified domain name where the record can be located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "record_hash",
				Description: "A unique identifier for the record within its zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "rdata",
				Description: "The record's data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Rdata"),
			},
			{
				Name:        "is_protected",
				Description: "A Boolean flag indicating whether or not parts of the record are unable to be explicitly managed.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "rrset_version",
				Description: "The latest version of the record's zone in which its RRSet differs from the preceding version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RrsetVersion"),
			},
			{
				Name:        "rtype",
				Description: "The type of DNS record, such as A or CNAME.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Rtype"),
			},
			{
				Name:        "ttl",
				Description: "The Time To Live for the record, in seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ttl"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Domain"),
			},

			// OCI standard columns
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

type dnsRecordInfo struct {
	dns.Record
	CompartmentId string
}

//// LIST FUNCTION

func listDnsRecordSets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("dns.listDnsRecordSets", "Compartment", compartment)

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

	zone := h.Item.(dns.ZoneSummary)

	request := dns.GetZoneRecordsRequest{
		ZoneNameOrId:  zone.Id,
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DnsClient.GetZoneRecords(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, record := range response.Items {
			d.StreamListItem(ctx, dnsRecordInfo{record, compartment})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
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
