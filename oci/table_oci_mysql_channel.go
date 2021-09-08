package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableMySQLChannel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_channel",
		Description: "OCI MySQL Channel",
		List: &plugin.ListConfig{
			Hydrate: listMySQLChannels,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getMySQLChannel,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The user-friendly name for the Channel. It does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the Channel.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Channel.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_enabled",
				Description: "Whether the Channel has been enabled by the user.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the Channel was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMySQLChannel,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "A user-supplied description of the backup.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLChannel,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the state of the Channel.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLChannel,
			},
			{
				Name:        "time_updated",
				Description: "The time the Channel was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMySQLChannel,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "source",
				Description: "Parameters detailing how to provision the source for the given Channel.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "target",
				Description: "Parameters detailing how to provision the target for the given Channel.",
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

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(channelTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// OCI standard columns
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
				Hydrate:     getMySQLChannel,
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

func listMySQLChannels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listMySQLChannels", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := mySQLChannelService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.ListChannelsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.MySQLChannelClient.ListChannels(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, channel := range response.Items {
			d.StreamListItem(ctx, channel)

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

//// HYDRATE FUNCTIONS

func getMySQLChannel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getMySQLChannel", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(mysql.ChannelSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := mySQLChannelService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.GetChannelRequest{
		ChannelId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.MySQLChannelClient.GetChannel(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Channel, nil
}

//// TRANSFORM FUNCTION

func channelTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := channelFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := channelDefinedTags(d.HydrateItem)

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

func channelFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case mysql.Channel:
		return item.FreeformTags
	case mysql.ChannelSummary:
		return item.FreeformTags
	}
	return nil
}

func channelDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case mysql.Channel:
		return item.DefinedTags
	case mysql.ChannelSummary:
		return item.DefinedTags
	}
	return nil
}
