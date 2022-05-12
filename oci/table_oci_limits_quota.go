package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/limits"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableLimitsQuota(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_limits_quota",
		Description: "OCI Limits Quota",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			Hydrate:           getLimitsQuota,
			ShouldIgnoreError: isNotFoundError([]string{"400", "404", "InvalidParameter"}),
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"InvalidParameter"}),
			Hydrate:           listLimitsQuotas,
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
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name you assign to the quota during creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.Name"),
			},
			{
				Name:        "id",
				Description: "The OCID of the quota.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.Id"),
			},
			{
				Name:        "description",
				Description: "The description you assign to the quota.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.Description"),
			},
			{
				Name:        "time_created",
				Description: "Date and time the quota was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Quota.TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The quota's current state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.LifecycleState"),
			},
			{
				Name:        "statements",
				Description: "An array of one or more quota statements written in the declarative quota statement language.",
				Hydrate:     getLimitsQuota,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Quota.Statements"),
			},

			// Tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Quota.DefinedTags"),
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Quota.FreeformTags"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(quotaTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.Name"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Quota.CompartmentId"),
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

type LimitsQuota struct {
	Quota  limits.Quota
	Region string
}

func listLimitsQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// fetch home region for the account
	gethomeRegionCached := plugin.HydrateFunc(getHomeRegion).WithCache()
	regions, err := gethomeRegionCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	region := regions.(string)
	if region == "" {
		return nil, nil
	}

	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listLimitsQuotas", "Compartment", compartment)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := quotaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := limits.ListQuotasRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	if d.KeyColumnQualString("name") != "" {
		name := d.KeyColumnQualString("name")
		request.Name = &name
	}
	if d.KeyColumnQualString("lifecycle_state") != "" {
		lifecycleState := d.KeyColumnQualString("lifecycle_state")
		request.LifecycleState = limits.ListQuotasLifecycleStateEnum(lifecycleState)
	}

	// Check for limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}
	pagesLeft := true
	for pagesLeft {

		response, err := session.QuotaClient.ListQuotas(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, quota := range response.Items {
			d.StreamListItem(ctx, &LimitsQuota{
				Region: region,
				Quota: limits.Quota{
					Id:             quota.Id,
					CompartmentId:  quota.CompartmentId,
					Name:           quota.Name,
					Description:    quota.Description,
					TimeCreated:    quota.TimeCreated,
					LifecycleState: limits.QuotaLifecycleStateEnum(quota.LifecycleState),
					FreeformTags:   quota.FreeformTags,
					DefinedTags:    quota.DefinedTags,
				},
			})

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

func getLimitsQuota(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getLimitsQuota")
	logger := plugin.Logger(ctx)

	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getLimitsQuota", "Compartment", compartment)

	// fetch home region for the account
	gethomeRegionCached := plugin.HydrateFunc(getHomeRegion).WithCache()
	regions, err := gethomeRegionCached(ctx, d, h)
	region := regions.(string)
	if err != nil {
		return nil, err
	}
	if region == "" {
		return nil, nil
	}

	// Create Session
	session, err := quotaService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	var id string
	if h.Item != nil {
		id = *h.Item.(*LimitsQuota).Quota.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// Empty check
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	request := limits.GetQuotaRequest{
		QuotaId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.QuotaClient.GetQuota(ctx, request)
	if err != nil {
		return nil, err
	}

	return &LimitsQuota{Quota: response.Quota, Region: region}, nil
}

//// TRANSFORM FUNCTION

func quotaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	quota := d.HydrateItem.(*LimitsQuota)

	var tags map[string]interface{}

	if quota != nil {
		tags = map[string]interface{}{}
		for k, v := range quota.Quota.FreeformTags {
			tags[k] = v
		}

		for _, v := range quota.Quota.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
