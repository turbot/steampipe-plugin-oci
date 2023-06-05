package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/autoscaling"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// // TABLE DEFINITION
func tableAutoScalingPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_autoscaling_auto_scaling_policy",
		Description:      "OCI Auto Scaling Policy",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"id", "auto_scaling_configuration_id"}),
			Hydrate:    getAutoscalingAutoScalingPolicy,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAutoScalingConfigurations,
			Hydrate:       listAutoscalingAutoScalingPolicies,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "auto_scaling_configuration_id",
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
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The ID of the autoscaling policy that is assigned after creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auto_scaling_configuration_id",
				Description: "The OCID of the autoscaling configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_enabled",
				Description: "Whether the autoscaling policy is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "policy_type",
				Description: "The type of autoscaling policy.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Auto Scaling Policy was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "capacity",
				Description: "The capacity requirements of the autoscaling policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAutoscalingAutoScalingPolicy,
			},

			// Standard Steampipe columns
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

type autoscalingpolicyInfo struct {
	Id                         *string
	DisplayName                *string
	TimeCreated                *common.SDKTime
	Capacity                   *autoscaling.Capacity
	IsEnabled                  *bool
	PolicyType                 *string
	AutoScalingConfigurationId *string
}

//// LIST FUNCTION

func listAutoscalingAutoScalingPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	configuration := h.Item.(autoscaling.AutoScalingConfigurationSummary)
	region := d.EqualsQualString(matrixKeyRegion)
	logger.Debug("listAutoscalingAutoScalingPolicies", "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Create Session
	session, err := autoScalingService(ctx, d, region)
	if err != nil {
		logger.Error("oci_autoscaling_auto_scaling_policy.listAutoscalingAutoScalingPolicies", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildAutoscalingAutoScalingPolicyFilters(equalQuals)
	request.Limit = types.Int(100)
	request.AutoScalingConfigurationId = configuration.Id
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
		response, err := session.AutoScalingClient.ListAutoScalingPolicies(ctx, request)
		if err != nil {
			logger.Error("oci_autoscaling_auto_scaling_policy.listAutoscalingAutoScalingPolicies", "api_error", err)
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, &autoscalingpolicyInfo{
				Id:                         respItem.Id,
				DisplayName:                respItem.DisplayName,
				PolicyType:                 respItem.PolicyType,
				IsEnabled:                  respItem.IsEnabled,
				AutoScalingConfigurationId: configuration.Id,
			})

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

func getAutoscalingAutoScalingPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	logger.Debug("getAutoscalingAutoScalingPolicy", "OCI_REGION", region)

	var id, configurationId string
	if h.Item != nil {
		data := h.Item.(*autoscalingpolicyInfo)
		id = *data.Id
		configurationId = *data.AutoScalingConfigurationId
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		configurationId = d.EqualsQuals["auto_scaling_configuration_id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" || configurationId == "" {
		return nil, nil
	}

	// Create Session

	session, err := autoScalingService(ctx, d, region)
	if err != nil {
		logger.Error("getAutoscalingAutoScalingPolicy", "error_AutoscalingService", err)
		return nil, err
	}

	request := autoscaling.GetAutoScalingPolicyRequest{
		AutoScalingPolicyId:        types.String(id),
		AutoScalingConfigurationId: types.String(configurationId),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.AutoScalingClient.GetAutoScalingPolicy(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.AutoScalingPolicy, nil
}

// Build additional filters
func buildAutoscalingAutoScalingPolicyFilters(equalQuals plugin.KeyColumnEqualsQualMap) autoscaling.ListAutoScalingPoliciesRequest {
	request := autoscaling.ListAutoScalingPoliciesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	return request
}
