package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/aianomalydetection"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAiAnomalyDetectionAiPrivateEndpoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_ai_anomaly_detection_ai_private_endpoint",
		Description:      "OCI AI Anomaly Detection AI Private Endpoint",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAiAnomalyDetectionAiPrivateEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listAiAnomalyDetectionAiPrivateEndpoints,
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
					Name:    "display_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Unique identifier that is immutable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "Private Reverse Connection Endpoint display name.",
				Type:        proto.ColumnType_STRING,
			},
						{
				Name:        "time_created",
				Description: "Time that AI Private Endpoint was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "subnet_id",
				Description: "Subnet Identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_zones",
				Description: "List of DNS zones to be used by the data assets.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the private endpoint resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail. For example, can be used to provide actionable information for a resource in 'Failed' state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "freeform_tags",
				Description: "Simple key-value pair that is applied without any predefined name, type or scope. Exists for cross-compatibility only.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a namespace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attached_data_assets",
				Description: "The list of dataAssets using the private reverse connection endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAiAnomalyDetectionAiPrivateEndpoint,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(aiAnomalyDetectionAiPrivateEndpointTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listAiAnomalyDetectionAiPrivateEndpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_ai_anomaly_detection_ai_private_endpoint.listAiAnomalyDetectionAiPrivateEndpoints", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := aiAnomalyDetectionService(ctx, d, region)
	if err != nil {
		logger.Error("oci_ai_anomaly_detection_ai_private_endpoint.listAiAnomalyDetectionAiPrivateEndpoints", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildAiAnomalyDetectionAiPrivateEndpointFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(100)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.AnomalyDetectionClient.ListAiPrivateEndpoints(ctx, request)
		if err != nil {
			logger.Error("oci_ai_anomaly_detection_ai_private_endpoint.listAiAnomalyDetectionAiPrivateEndpoints", "api_error", err)
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

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

//// HYDRATE FUNCTIONS

func getAiAnomalyDetectionAiPrivateEndpoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_ai_anomaly_detection_ai_private_endpoint.getAiAnomalyDetectionAiPrivateEndpoint", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(aianomalydetection.AiPrivateEndpointSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := aiAnomalyDetectionService(ctx, d, region)
	if err != nil {
		logger.Error("oci_ai_anomaly_detection_ai_private_endpoint.getAiAnomalyDetectionAiPrivateEndpoint", "connection_error", err)
		return nil, err
	}

	request := aianomalydetection.GetAiPrivateEndpointRequest{
		AiPrivateEndpointId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AnomalyDetectionClient.GetAiPrivateEndpoint(ctx, request)
	if err != nil {
		logger.Error("oci_ai_anomaly_detection_ai_private_endpoint.getAiAnomalyDetectionAiPrivateEndpoint", "api_error", err)
		return nil, err
	}
	return response.AiPrivateEndpoint, nil
}

//// TRANSFORM FUNCTIONS

func aiAnomalyDetectionAiPrivateEndpointTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case aianomalydetection.AiPrivateEndpoint:
		obj := d.HydrateItem.(aianomalydetection.AiPrivateEndpoint)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case aianomalydetection.AiPrivateEndpointSummary:
		obj := d.HydrateItem.(aianomalydetection.AiPrivateEndpointSummary)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
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

// Build additional filters
func buildAiAnomalyDetectionAiPrivateEndpointFilters(equalQuals plugin.KeyColumnEqualsQualMap) aianomalydetection.ListAiPrivateEndpointsRequest {
	request := aianomalydetection.ListAiPrivateEndpointsRequest{}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = aianomalydetection.AiPrivateEndpointLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	return request
}
