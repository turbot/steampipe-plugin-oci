package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/autoscaling"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// // TABLE DEFINITION
func tableAutoscalingAutoScalingPolicy(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_autoscaling_auto_scaling_policy",
		Description:      "OCI Auto Scaling Policy",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAutoscalingAutoScalingPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listAutoscalingAutoScalingPolicies,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "capacity",
				Description: "The capacity requirements of the autoscaling policy.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAutoscalingAutoScalingPolicy,
			},
			{
				Name:        "id",
				Description: "The ID of the autoscaling policy that is assigned after creation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
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

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
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
func listAutoscalingAutoScalingPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listAutoscalingAutoScalingPolicies", "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	// Create Session
	session, err := autoscalingService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildAutoscalingAutoScalingPolicyFilters(equalQuals)
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
		response, err := session.AutoScalingClient.ListAutoScalingPolicies(ctx, request)
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
func getAutoscalingAutoScalingPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getAutoscalingAutoScalingPolicy", "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(autoscaling.AutoScalingPolicySummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := autoscalingService(ctx, d, region)
	if err != nil {
		logger.Error("getAutoscalingAutoScalingPolicy", "error_AutoscalingService", err)
		return nil, err
	}

	request := autoscaling.GetAutoScalingPolicyRequest{
		AutoScalingPolicyId: types.String(id),
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
