package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/analytics"
	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAnalyticsInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_analytics_instance",
		Description: "OCI Analytics Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAnalyticsInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsInstances,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "capacity_type",
					Require: plugin.Optional,
				},
				{
					Name:    "feature_set",
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
		GetMatrixItemFunc: BuildCompartmentRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The analytics instance's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The analytics instance's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the instance was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "capacity_type",
				Description: "The analytics instance's capacity model to use.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Capacity.CapacityType"),
			},
			{
				Name:        "description",
				Description: "The analytics instance's optional description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "email_notification",
				Description: "The email address receiving notifications.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "feature_set",
				Description: "The analytics instance's feature set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_type",
				Description: "The license used for the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_url",
				Description: "The URL of the Analytics service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "capacity_value",
				Description: "The analytics instance's capacity value selected.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Capacity.CapacityValue"),
			},
			{
				Name:        "network_endpoint_details",
				Description: "The base representation of a network endpoint.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_access_channels",
				Description: "The private access channels of the analytics instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAnalyticsInstance,
			},
			{
				Name:        "vanity_url_details",
				Description: "The vanity url configuration details of the analytic instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAnalyticsInstance,
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAnalyticsInstance,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAnalyticsInstance,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAnalyticsInstance,
				Transform:   transform.From(analyticsInstanceTags),
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

func listAnalyticsInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listAnalyticsInstances", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := analyticsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildAnalyticsInstanceFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
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
		response, err := session.AnalyticsClient.ListAnalyticsInstances(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, instance := range response.Items {
			d.StreamListItem(ctx, instance)

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

func getAnalyticsInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAnalyticsInstance")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getAnalyticsInstance", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(analytics.AnalyticsInstanceSummary).Id
	} else {
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}

		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty analytics instance id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := analyticsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := analytics.GetAnalyticsInstanceRequest{
		AnalyticsInstanceId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AnalyticsClient.GetAnalyticsInstance(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.AnalyticsInstance, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Free-form tags
// 2. Defined Tags
func analyticsInstanceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(analytics.AnalyticsInstance)

	var tags map[string]interface{}

	if instance.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range instance.FreeformTags {
			tags[k] = v
		}
	}

	if instance.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range instance.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

// Build additional filters
func buildAnalyticsInstanceFilters(equalQuals plugin.KeyColumnEqualsQualMap) analytics.ListAnalyticsInstancesRequest {
	request := analytics.ListAnalyticsInstancesRequest{}

	if equalQuals["capacity_type"] != nil {
		request.CapacityType = analytics.ListAnalyticsInstancesCapacityTypeEnum(equalQuals["capacity_type"].GetStringValue())
	}
	if equalQuals["feature_set"] != nil {
		request.FeatureSet = analytics.ListAnalyticsInstancesFeatureSetEnum(equalQuals["feature_set"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = analytics.ListAnalyticsInstancesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}
	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}

	return request
}
