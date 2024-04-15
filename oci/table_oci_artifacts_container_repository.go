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

//// TABLE DEFINITION

func tableArtifactContainerRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifacts_container_repository",
		Description:      "OCI Artifacts Container Repository",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getArtifactContainerRepository,
		},
		List: &plugin.ListConfig{
			Hydrate: listArtifactContainerRepositories,
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
					Name:    "is_public",
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
				Name:        "display_name",
				Description: "The container repository name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the container repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "The id of the user or principal that created the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactContainerRepository,
			},
			{
				Name:        "time_created",
				Description: "Time that Container Repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "image_count",
				Description: "Total number of images.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "is_immutable",
				Description: "Whether the repository is immutable. Images cannot be overwritten in an immutable repository.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getArtifactContainerRepository,
			},
			{
				Name:        "is_public",
				Description: "Whether the repository is public. A public repository allows unauthenticated access.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "layer_count",
				Description: "Total number of layers.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "layers_size_in_bytes",
				Description: "Total storage in bytes consumed by layers.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the container repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "billable_size_in_gbs",
				Description: "Total storage size in GBs that will be charged.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BillableSizeInGBs"),
			},
			{
				Name:        "time_last_pushed",
				Description: "An RFC 3339 timestamp indicating when an image was last pushed to the repository.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getArtifactContainerRepository,
				Transform:   transform.FromField("TimeLastPushed.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,

				Transform: transform.FromField("DisplayName"),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listArtifactContainerRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifacts_container_repository.listArtifactContainerRepositories", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_container_repository.listArtifactContainerRepositories", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildArtifactContainerRepositoryFilters(equalQuals)
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
		response, err := session.ArtifactClient.ListContainerRepositories(ctx, request)
		if err != nil {
			logger.Error("oci_artifacts_container_repository.listArtifactContainerRepositories", "api_error", err)
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

func getArtifactContainerRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifacts_container_repository.getArtifactContainerRepository", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(artifacts.ContainerRepositorySummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty ID in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_container_repository.getArtifactContainerRepository", "connection_error", err)
		return nil, err
	}

	request := artifacts.GetContainerRepositoryRequest{
		RepositoryId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ArtifactClient.GetContainerRepository(ctx, request)
	if err != nil {
		logger.Error("oci_artifacts_container_repository.getArtifactContainerRepository", "api_error", err)
		return nil, err
	}
	return response.ContainerRepository, nil
}

// Build additional filters
func buildArtifactContainerRepositoryFilters(equalQuals plugin.KeyColumnEqualsQualMap) artifacts.ListContainerRepositoriesRequest {
	request := artifacts.ListContainerRepositoriesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["is_public"] != nil {
		request.IsPublic = types.Bool(equalQuals["is_public"].GetBoolValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = types.String(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
