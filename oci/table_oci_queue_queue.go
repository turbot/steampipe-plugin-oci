package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/queue"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableQueueQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_queue_queue",
		Description: "OCI Queue Queue",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getQueue,
		},
		List: &plugin.ListConfig{
			Hydrate: listQueues,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The name of the queue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "The OCID of the queue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "messages_endpoint",
				Description: "The endpoint to use to get or put messages in the queue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail. For example, can be used to provide actionable information for a resource in Failed state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Queue.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the queue was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the queue was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "retention_in_seconds",
				Description: "The retention period of the messages in the queue, in seconds.",
				Hydrate:     getQueue,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "visibility_in_seconds",
				Description: "The default visibility of the messages consumed from the queue.",
				Hydrate:     getQueue,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "timeout_in_seconds",
				Description: "The default polling timeout of the messages in the queue, in seconds.",
				Hydrate:     getQueue,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "dead_letter_queue_delivery_count",
				Description: "The number of times a message can be delivered to a consumer before being moved to the dead letter queue. A value of 0 indicates that the DLQ is not used.",
				Hydrate:     getQueue,
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "custom_encryption_key_id ",
				Description: "ID of the custom master encryption key which will be used to encrypt messages content.",
				Hydrate:     getQueue,
				Type:        proto.ColumnType_INT,
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

			// // Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(queueTags),
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

func listQueues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)

	logger.Debug("listQueues", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := queueService(ctx, d, region)
	if err != nil {
                logger.Error("oci_queue_queue. listQueues", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request, isValid := buildQueueFilters(equalQuals, logger)
	if !isValid {
		return nil, nil
	}
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = oci_common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.QueueAdminClient.ListQueues(ctx, request)
		if err != nil {
                         logger.Error("oci_queue_queue. listQueues", "api_error", err)
			return nil, err
		}

		for _, queueSummary := range response.Items {
			d.StreamListItem(ctx, queueSummary)

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

//// HYDRATE FUNCTIONS

func getQueue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getQueue")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getQueue", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(queue.QueueSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	if id == "" {
		return nil, nil
	}
	// Create Session
	session, err := queueService(ctx, d, region)
	if err != nil {
                logger.Error("oci_queue_queue. getQueue", "connection_error", err)
		return nil, err
	}

	request := queue.GetQueueRequest{
		QueueId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.QueueAdminClient.GetQueue(ctx, request)
	if err != nil {
                logger.Error("oci_queue_queue. getQueue", "api_error", err)
		return nil, err
	}

	return response.Queue, nil
}

// Build additional filters
func buildQueueFilters(equalQuals plugin.KeyColumnEqualsQualMap, logger hclog.Logger) (queue.ListQueuesRequest, bool) {
	request := queue.ListQueuesRequest{}
	isValid := true

	if equalQuals["displayName"] != nil && strings.Trim(equalQuals["displayName"].GetStringValue(), " ") != "" {
		request.DisplayName = types.String(equalQuals["displayName"].GetStringValue())
	}
	return request, isValid
}

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func queueTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case queue.QueueSummary:
		freeFormTags = d.HydrateItem.(queue.QueueSummary).FreeformTags
		definedTags = d.HydrateItem.(queue.QueueSummary).DefinedTags
	case queue.Queue:
		freeFormTags = d.HydrateItem.(queue.Queue).FreeformTags
		definedTags = d.HydrateItem.(queue.Queue).DefinedTags
	default:
		return nil, nil
	}

	var tags map[string]interface{}
	if freeFormTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeFormTags {
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
