package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/autoscaling"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strings"
)

// // TABLE DEFINITION
func tableAutoscalingAutoScalingConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_autoscaling_auto_scaling_configuration",
		Description:      "OCI Auto Scaling Configuration",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAutoscalingAutoScalingConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listAutoscalingAutoScalingConfigurations,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
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
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the autoscaling configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource",
				Description: "A resource that is managed by an autoscaling configuration. The only supported type is instancePool.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "policies",
				Description: "Autoscaling policy definitions for the autoscaling configuration. An autoscaling policy defines the criteria that",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAutoscalingAutoScalingConfiguration,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource. Each tag is a simple key-value pair with no",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "cool_down_in_seconds",
				Description: "For threshold-based autoscaling policies, this value is the minimum period of time to wait between scaling actions.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_enabled",
				Description: "Whether the autoscaling configuration is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "max_resource_count",
				Description: "The maximum number of resources to scale out to.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAutoscalingAutoScalingConfiguration,
			},
			{
				Name:        "min_resource_count",
				Description: "The minimum number of resources to scale in to.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAutoscalingAutoScalingConfiguration,
			},
			{
				Name:        "time_created",
				Description: "Time that Auto Scaling Configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(autoscalingAutoScalingConfigurationTags),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

// // LIST FUNCTION
func listAutoscalingAutoScalingConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listAutoscalingAutoScalingConfigurations", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := autoscalingService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildAutoscalingAutoScalingConfigurationFilters(equalQuals)
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
		response, err := session.AutoScalingClient.ListAutoScalingConfigurations(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

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

// // HYDRATE FUNCTION
func getAutoscalingAutoScalingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getAutoscalingAutoScalingConfiguration", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(autoscaling.AutoScalingConfigurationSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := autoscalingService(ctx, d, region)
	if err != nil {
		logger.Error("getAutoscalingAutoScalingConfiguration", "error_AutoscalingService", err)
		return nil, err
	}

	request := autoscaling.GetAutoScalingConfigurationRequest{
		AutoScalingConfigurationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AutoScalingClient.GetAutoScalingConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.AutoScalingConfiguration, nil
}

// // TRANSFORM FUNCTION
func autoscalingAutoScalingConfigurationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case autoscaling.AutoScalingConfiguration:
		obj := d.HydrateItem.(autoscaling.AutoScalingConfiguration)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case autoscaling.AutoScalingConfigurationSummary:
		obj := d.HydrateItem.(autoscaling.AutoScalingConfigurationSummary)
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
func buildAutoscalingAutoScalingConfigurationFilters(equalQuals plugin.KeyColumnEqualsQualMap) autoscaling.ListAutoScalingConfigurationsRequest {
	request := autoscaling.ListAutoScalingConfigurationsRequest{}

	if equalQuals["compartment_id"] != nil {
		request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	return request
}
