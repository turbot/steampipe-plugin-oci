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

func tableDnsTsigKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_dns_tsig_key",
		Description: "OCI DNS TSIG Key",
		List: &plugin.ListConfig{
			Hydrate: listDnsTsigKeys,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDnsTsigKey,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A globally unique domain name identifying the key for a given pair of hosts.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "algorithm",
				Description: "TSIG key algorithms are encoded as domain names, but most consist of only one non-empty label, which is not required to be explicitly absolute.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "secret",
				Description: "A base64 string encoding the binary shared secret.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDnsTsigKey,
			},
			{
				Name:        "self",
				Description: "The canonical absolute URL of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The date and time the resource was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
				Hydrate:     getDnsTsigKey,
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
				Transform:   transform.From(dnsTsigTags),
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
				Transform:   transform.FromField("Self").Transform(ociRegionFromSelf),
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

func listDnsTsigKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listDnsTsigKeys", "Compartment", compartment)

	// Create Session
	session, err := dnsService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := dns.ListTsigKeysRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		keys, err := session.DnsClient.ListTsigKeys(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, key := range keys.Items {
			d.StreamListItem(ctx, key)
		}
		if keys.OpcNextPage != nil {
			request.Page = keys.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDnsTsigKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getDnsTsigKey", "Compartment", compartment)

	// Rstrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	var id string
	if h.Item != nil {
		id = *h.Item.(dns.TsigKeySummary).Id
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

	request := dns.GetTsigKeyRequest{
		TsigKeyId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.DnsClient.GetTsigKey(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.TsigKey, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 2. Defined Tags
// 3. Free-form tags
func dnsTsigTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case dns.TsigKey:
		key := d.HydrateItem.(dns.TsigKey)
		freeformTags = key.FreeformTags
		definedTags = key.DefinedTags
	case dns.TsigKeySummary:
		key := d.HydrateItem.(dns.TsigKeySummary)
		freeformTags = key.FreeformTags
		definedTags = key.DefinedTags
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