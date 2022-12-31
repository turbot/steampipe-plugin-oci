package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/artifacts"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strings"
)

// // TABLE DEFINITION
func tableArtifactsRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifacts_repository",
		Description:      "OCI Repository",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getArtifactsRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listArtifactsRepositories,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "is_immutable",
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
				Description: "TBC",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "display_name",
				Description: "TBC",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "description",
				Description: "TBC",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "is_immutable",
				Description: "TBC",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "lifecycle_state",
				Description: "TBC",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "freeform_tags",
				Description: "TBC",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "defined_tags",
				Description: "TBC",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "repository_type",
				Description: "TBC",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactsRepository,
			},
			{
				Name:        "time_created",
				Description: "Time that Repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(artifactsRepositoryTags),
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
func listArtifactsRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listArtifactsRepositories", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := artifactsService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildArtifactsRepositoryFilters(equalQuals)
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
		response, err := session.ArtifactsClient.ListRepositories(ctx, request)
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
func getArtifactsRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getArtifactsRepository", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(artifacts.RepositorySummary).GetId()
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

	session, err := artifactsService(ctx, d, region)
	if err != nil {
		logger.Error("getArtifactsRepository", "error_ArtifactsService", err)
		return nil, err
	}

	request := artifacts.GetRepositoryRequest{
		RepositoryId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ArtifactsClient.GetRepository(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Repository, nil
}

// // TRANSFORM FUNCTION
func artifactsRepositoryTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case artifacts.Repository:
		obj := d.HydrateItem.(artifacts.Repository)
		freeformTags = obj.GetFreeformTags()
		definedTags = obj.GetDefinedTags()
	case artifacts.RepositorySummary:
		obj := d.HydrateItem.(artifacts.RepositorySummary)
		freeformTags = obj.GetFreeformTags()
		definedTags = obj.GetDefinedTags()
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
func buildArtifactsRepositoryFilters(equalQuals plugin.KeyColumnEqualsQualMap) artifacts.ListRepositoriesRequest {
	request := artifacts.ListRepositoriesRequest{}

	if equalQuals["compartment_id"] != nil {
		request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())
	}

	if equalQuals["id"] != nil {
		request.Id = types.String(equalQuals["id"].GetStringValue())
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["is_immutable"] != nil {
		request.IsImmutable = types.Bool(equalQuals["is_immutable"].GetBoolValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = types.String(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
