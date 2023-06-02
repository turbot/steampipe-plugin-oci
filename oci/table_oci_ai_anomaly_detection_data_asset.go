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

func tableAiAnomalyDetectionDataAsset(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_ai_anomaly_detection_data_asset",
		Description:      "OCI Anomaly Detection Data Asset",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAiAnomalyDetectionDataAsset,
		},
		List: &plugin.ListConfig{
			Hydrate: listAiAnomalyDetectionDataAssets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "project_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
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
				Name:        "id",
				Description: "The Unique Oracle ID (OCID) that is immutable on creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The lifecycle state of the Data Asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Data Asset was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "project_id",
				Description: "The Unique project id which is created at project creation that is immutable on creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A short description of the data asset.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_endpoint_id",
				Description: "OCID of Private Endpoint.",
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(aiAnomalyDetectionDataAssetTags),
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

func listAiAnomalyDetectionDataAssets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_ai_anomaly_detection_data_asset.listAiAnomalyDetectionDataAssets", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := aiAnomalyDetectionService(ctx, d, region)
	if err != nil {
		logger.Error("oci_ai_anomaly_detection_data_asset.listAiAnomalyDetectionDataAssets", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildAiAnomalyDetectionDataAssetFilters(equalQuals)
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
		response, err := session.AnomalyDetectionClient.ListDataAssets(ctx, request)
		if err != nil {
			logger.Error("oci_ai_anomaly_detection_data_asset.listAiAnomalyDetectionDataAssets", "api_error", err)
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

func getAiAnomalyDetectionDataAsset(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_ai_anomaly_detection_data_asset.getAiAnomalyDetectionDataAsset", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(aianomalydetection.DataAssetSummary).Id
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
		logger.Error("oci_ai_anomaly_detection_data_asset.getAiAnomalyDetectionDataAsset", "connection_error", err)
		return nil, err
	}

	request := aianomalydetection.GetDataAssetRequest{
		DataAssetId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AnomalyDetectionClient.GetDataAsset(ctx, request)
	if err != nil {
		logger.Error("oci_ai_anomaly_detection_data_asset.getAiAnomalyDetectionDataAsset", "api_error", err)
		return nil, err
	}
	return response.DataAsset, nil
}

//// TRANSFORM FUNCTIONS

func aiAnomalyDetectionDataAssetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case aianomalydetection.DataAsset:
		obj := d.HydrateItem.(aianomalydetection.DataAsset)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case aianomalydetection.DataAssetSummary:
		obj := d.HydrateItem.(aianomalydetection.DataAssetSummary)
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
func buildAiAnomalyDetectionDataAssetFilters(equalQuals plugin.KeyColumnEqualsQualMap) aianomalydetection.ListDataAssetsRequest {
	request := aianomalydetection.ListDataAssetsRequest{}

	if equalQuals["project_id"] != nil {
		request.ProjectId = types.String(equalQuals["project_id"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = aianomalydetection.DataAssetLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
