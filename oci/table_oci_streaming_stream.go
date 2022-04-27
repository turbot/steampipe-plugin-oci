package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/streaming"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableOciStreamingStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_streaming_stream",
		Description: "OCI Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getStreamingStream,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listStreamingStreams,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "stream_pool_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the stream.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the stream was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_state_details",
				Description: "Any additional details about the current state of the stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStreamingStream,
			},
			{
				Name:        "messages_endpoint",
				Description: "The endpoint to use when creating the StreamClient to consume or publish messages in the stream.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getStreamingStream,
			},
			{
				Name:        "partitions",
				Description: "The number of partitions in the stream.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "retention_in_hours",
				Description: "The retention period of the stream, in hours. This property is read-only.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getStreamingStream,
			},
			{
				Name:        "stream_pool_id",
				Description: "The OCID of the stream pool that contains the stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("StreamPoolId"),
			},

			// Tags
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
				Transform:   transform.From(streamTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

func listStreamingStreams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listStreams", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := streamAdminService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := streaming.ListStreamsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(50),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["name"] != nil {
		name := equalQuals["name"].GetStringValue()
		request.Name = types.String(name)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidStreamLifecycleStateEnum(lifecycleState) {
			request.LifecycleState = streaming.StreamLifecycleStateEnum(lifecycleState)
		} else {
			return nil, nil
		}
	}

	if equalQuals["stream_pool_id"] != nil {
		streamPoolId := equalQuals["stream_pool_id"].GetStringValue()
		request.StreamPoolId = types.String(streamPoolId)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.StreamAdminClient.ListStreams(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, stream := range response.Items {
			d.StreamListItem(ctx, stream)

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

func getStreamingStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getStreams", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(streaming.StreamSummary).Id
	} else {

		// Restrict the api call to only root compartment
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := streamAdminService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := streaming.GetStreamRequest{
		StreamId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.StreamAdminClient.GetStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Stream, nil
}

//// TRANSFORM FUNCTION

func streamTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case streaming.StreamSummary:
		streamSummary := d.HydrateItem.(streaming.StreamSummary)
		freeformTags = streamSummary.FreeformTags
		definedTags = streamSummary.DefinedTags
	case streaming.Stream:
		stream := d.HydrateItem.(streaming.Stream)
		freeformTags = stream.FreeformTags
		definedTags = stream.DefinedTags
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

func isValidStreamLifecycleStateEnum(state string) bool {
	stateType := streaming.StreamLifecycleStateEnum(state)

	switch stateType {
	case streaming.StreamLifecycleStateCreating, streaming.StreamLifecycleStateActive, streaming.StreamLifecycleStateDeleted, streaming.StreamLifecycleStateDeleting, streaming.StreamLifecycleStateFailed, streaming.StreamLifecycleStateUpdating:
		return true
	}
	return false
}
