package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/artifacts"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// // TABLE DEFINITION
func tableArtifactsRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifacts_repository",
		Description:      "OCI Artifact Repository",
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
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The repository name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "An RFC 3339 timestamp indicating when the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "The repository description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_immutable",
				Description: "Whether the repository is immutable. The artifacts of an immutable repository cannot be overwritten.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource. Each tag is a simple key-value pair with no",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a",
				Type:        proto.ColumnType_JSON,
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
				Description: ColumnDescriptionTenantId,
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
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifacts_repository.listArtifactsRepositories", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := artifactsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_repository.listArtifactsRepositories", "connection_error", err)
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
			logger.Error("oci_artifacts_repository.listArtifactsRepositories", "api_error", err)
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

//// HYDRATE FUNCTION

func getArtifactsRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getArtifactsRepository", "Compartment", compartment, "OCI_REGION", region)

	id := d.EqualsQuals["id"].GetStringValue()

	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := artifactsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_repository.getArtifactsRepository", "connection_error", err)
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
		logger.Error("oci_artifacts_repository.getArtifactsRepository", "api_error", err)
		return nil, err
	}
	return response.Repository, nil
}

// // TRANSFORM FUNCTION
func artifactsRepositoryTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case artifacts.RepositorySummary:
		obj := d.HydrateItem.(artifacts.RepositorySummary)
		freeformTags = obj.GetFreeformTags()
		definedTags = obj.GetDefinedTags()
	case artifacts.Repository:
		obj := d.HydrateItem.(artifacts.Repository)
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
