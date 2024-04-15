package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/ons"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOnsNotificationTopic(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_ons_notification_topic",
		Description: "A communication channel for sending messages to the subscriptions.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("topic_id"),
			Hydrate:    getOnsNotificationTopic,
		},
		List: &plugin.ListConfig{
			Hydrate: listOnsNotificationTopics,
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
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "topic_id",
				Description: "The OCID of the topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The lifecycle state of the topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the topic was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "api_endpoint",
				Description: "The endpoint for managing subscriptions or publishing messages to the topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "The description of the topic.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "Used for optimistic concurrency control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "short_topic_id",
				Description: "A unique short topic Id. This is used only for SMS subscriptions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
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
				Transform:   transform.From(topicTags),
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
				Transform:   transform.FromField("TopicId").Transform(ociRegionName),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listOnsNotificationTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.listOnsNotificationTopics", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := onsNotificationControlPlaneService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := ons.ListTopicsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(50),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = ons.NotificationTopicSummaryLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.NotificationControlPlaneClient.ListTopics(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, topic := range response.Items {
			d.StreamListItem(ctx, topic)

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

func getOnsNotificationTopic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOnsNotificationTopic")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.getOnsNotificationTopic", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQuals["topic_id"].GetStringValue()

	// handle empty topic id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := onsNotificationControlPlaneService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := ons.GetTopicRequest{
		TopicId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.NotificationControlPlaneClient.GetTopic(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.NotificationTopic, nil
}

//// TRANSFORM FUNCTION

func topicTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := topicFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := topicDefinedTags(d.HydrateItem)

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

func topicFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case ons.NotificationTopic:
		return item.FreeformTags
	case ons.NotificationTopicSummary:
		return item.FreeformTags
	}
	return nil
}

func topicDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case ons.NotificationTopic:
		return item.DefinedTags
	case ons.NotificationTopicSummary:
		return item.DefinedTags
	}
	return nil
}
