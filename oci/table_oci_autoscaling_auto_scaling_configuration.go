package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/autoscaling"
	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAutoScalingConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_autoscaling_auto_scaling_configuration",
		Description: "OCI Auto Scaling Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAutoScalingConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listAutoScalingConfigurations,
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
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the autoscaling configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_enabled",
				Description: "Indicates whether the autoscaling configuration is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "time_created",
				Description: "The date and time the AutoScalingConfiguration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "cool_down_in_seconds",
				Description: "The minimum period of time to wait between scaling actions.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_resource_count",
				Description: "The maximum number of resources to scale out to.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAutoScalingConfiguration,
			},
			{
				Name:        "min_resource_count",
				Description: "The minimum number of resources to scale in to.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAutoScalingConfiguration,
			},

			// Json fields
			{
				Name:        "policies",
				Description: "Autoscaling policy definitions for the autoscaling configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAutoScalingConfiguration,
			},
			{
				Name:        "resource",
				Description: "The resource details of this AutoScalingConfiguration.",
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(autoScalingConfigurationTags),
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

func listAutoScalingConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listAutoScalingConfigurations", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := autoScalingService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := autoscaling.ListAutoScalingConfigurationsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Check for additional filter
	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
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

		for _, configuration := range response.Items {
			d.StreamListItem(ctx, configuration)

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

func getAutoScalingConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAutoScalingConfiguration")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getAutoScalingConfiguration", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	var id string
	if h.Item != nil {
		autoScalingConfiguration := h.Item.(autoscaling.AutoScalingConfigurationSummary)
		id = *autoScalingConfiguration.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty volume backup id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := autoScalingService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := autoscaling.GetAutoScalingConfigurationRequest{
		AutoScalingConfigurationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.AutoScalingClient.GetAutoScalingConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.AutoScalingConfiguration, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func autoScalingConfigurationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := autoScalingConfigurationFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := autoScalingConfigurationDefinedTags(d.HydrateItem)

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

func autoScalingConfigurationFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case autoscaling.AutoScalingConfiguration:
		return item.FreeformTags
	case autoscaling.AutoScalingConfigurationSummary:
		return item.FreeformTags
	}
	return nil
}

func autoScalingConfigurationDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case autoscaling.AutoScalingConfiguration:
		return item.DefinedTags
	case autoscaling.AutoScalingConfigurationSummary:
		return item.DefinedTags
	}
	return nil
}
