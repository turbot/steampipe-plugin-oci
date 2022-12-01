package oci

import (
	"context"
	"strings"


	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/ons"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableOnsSubscription(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_ons_subscription",
		Description: "OCI Ons Subscription",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"400", "404"}),
			Hydrate:           getOnsSubscription,
		},
		List: &plugin.ListConfig{
			Hydrate: listOnsSubscriptions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "topic_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartmentRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "topic_id",
				Description: "The OCID of the associated topic.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The lifecycle state of the subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time when this subscription was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromCamel().Transform(transform.UnixMsToTimestamp),
			},

			// other columns
			{
				Name:        "endpoint",
				Description: "A locator that corresponds to the subscription protocol.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "Used for optimistic concurrency control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol",
				Description: "The protocol used for the subscription.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "delivery_policy",
				Description: "Delivery Policy of the subscription.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSubscriptionDeliveryPolicy,
				Transform:   transform.FromValue(),
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
				Transform:   transform.From(subscriptionTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Endpoint"),
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

func listOnsSubscriptions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listOnsSubscriptions", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := onsNotificationDataPlaneService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := ons.ListSubscriptionsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(50),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filter
	if equalQuals["topic_id"] != nil {
		topicId := equalQuals["topic_id"].GetStringValue()
		request.TopicId = types.String(topicId)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.NotificationDataPlaneClient.ListSubscriptions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, subscription := range response.Items {
			d.StreamListItem(ctx, subscription)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func getOnsSubscription(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getOnsSubscription")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getOnsSubscription", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty subscription id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := onsNotificationDataPlaneService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := ons.GetSubscriptionRequest{
		SubscriptionId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.NotificationDataPlaneClient.GetSubscription(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Subscription, nil
}

/*
The delivery policy is returned in the list, but not get the call (https://github.com/turbot/steampipe-plugin-oci/issues/369).
So if we're hydrated from the list call, return the delivery policy directly,
but if we're hydrated from the get call, we need to make an extra list call and
filter on the topic ID.
*/
func getSubscriptionDeliveryPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getSubscriptionDeliveryPolicy")
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getOnsSubscription", "Compartment", compartment, "OCI_REGION", region)

	policy := deliveryPolicy(ctx, h.Item)

	// If the subscription is from the list call, we already have the policy
	if policy != nil {
		return policy, nil
	}

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := onsNotificationDataPlaneService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	subscriptionItem := h.Item.(ons.Subscription)
	request := ons.ListSubscriptionsRequest{
		CompartmentId: types.String(compartment),
		TopicId:       types.String(*subscriptionItem.TopicId),
		Limit:         types.Int(50),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.NotificationDataPlaneClient.ListSubscriptions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, subscription := range response.Items {
			if (*subscription.Id == (*subscriptionItem.Id)){
				return subscription.DeliveryPolicy, nil
			}
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	// We shouldn't hit this error condition unless the subscription is deleted
	// during the API call
	logger.Error("oci.getOnsSubscription", "subscription_not_found_error", err)
	return nil, err
}


//// TRANSFORM FUNCTION

func subscriptionTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := subscriptionFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := subscriptionDefinedTags(d.HydrateItem)

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

func subscriptionFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case ons.Subscription:
		return item.FreeformTags
	case ons.SubscriptionSummary:
		return item.FreeformTags
	}
	return nil
}

func subscriptionDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case ons.Subscription:
		return item.DefinedTags
	case ons.SubscriptionSummary:
		return item.DefinedTags
	}
	return nil
}

func deliveryPolicy(ctx context.Context, item interface{}) *ons.DeliveryPolicy {
	switch item := item.(type) {
		case ons.SubscriptionSummary:
		return item.DeliveryPolicy
	}
	return nil
}
