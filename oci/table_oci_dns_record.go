package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/dns"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableDnsRecord(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_dns_record",
		Description: "OCI Dns Record",
		List: &plugin.ListConfig{
			ParentHydrate: listDnsZones,
			Hydrate:       listDnsRecords,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "domain",
				Description: " The fully qualified domain name where the record can be located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "record_hash",
				Description: "A unique identifier for the record within its zone.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "r_data",
				Description: "The record's data.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Rdata"),
			},
			{
				Name:        "rrset_version",
				Description: "The latest version of the record's zone in which its RRSet differs from the preceding version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RrsetVersion"),
			},
			{
				Name:        "rtype",
				Description: "The date and time the zone was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Rtype"),
			},
			{
				Name:        "ttl",
				Description: "The Time To Live for the record, in seconds.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ttl"),
			},
			{
				Name:        "is_protected",
				Description: "A Boolean flag indicating whether or not parts of the record are unable to be explicitly managed.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Domain"),
			},

			// Standard OCI columns
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

func listDnsRecords(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("dns.listDnsRecords", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := dnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	zone := h.Item.(dns.ZoneSummary)

	request := dns.GetZoneRecordsRequest{
		ZoneNameOrId:  zone.Id,
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DnsClient.GetZoneRecords(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, zone := range response.Items {
			d.StreamLeafListItem(ctx, zone)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}
