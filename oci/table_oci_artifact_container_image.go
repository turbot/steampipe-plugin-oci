package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/artifacts"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

// // TABLE DEFINITION
func tableArtifactContainerImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifact_container_image",
		Description:      "OCI Container Image",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getArtifactContainerImage,
		},
		List: &plugin.ListConfig{
			Hydrate: listArtifactContainerImages,
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
					Name:    "repository_id",
					Require: plugin.Optional,
				},
				{
					Name:    "repository_name",
					Require: plugin.Optional,
				},
				{
					Name:    "version",
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
				Name:        "created_by",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the user or principal that created the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "digest",
				Description: "The container image digest.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The repository name and the most recent version associated with the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the container image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "layers",
				Description: "Layers of which the image is composed, ordered by the layer digest.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "layers_size_in_bytes",
				Description: "The total size of the container image layers in bytes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the container image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "manifest_size_in_bytes",
				Description: "The size of the container image manifest in bytes.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "pull_count",
				Description: "Total number of pulls.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "repository_id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the container repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_name",
				Description: "The container repository name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "versions",
				Description: "The versions associated with this image.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getArtifactContainerImage,
			},
			{
				Name:        "time_last_pulled",
				Description: "An RFC 3339 timestamp indicating when the image was last pulled.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getArtifactContainerImage,
				Transform:   transform.FromField("TimeLastPulled.Time"),
			},
			{
				Name:        "version",
				Description: "The most recent version associated with this image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Container Image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

// // LIST FUNCTION
func listArtifactContainerImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listArtifactContainerImages", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := artifactService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildArtifactContainerImageFilters(equalQuals)
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
		response, err := session.ArtifactClient.ListContainerImages(ctx, request)
		if err != nil {
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

// // HYDRATE FUNCTION
func getArtifactContainerImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getArtifactContainerImage", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(artifacts.ContainerImageSummary).Id
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

	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("getArtifactContainerImage", "error_ArtifactService", err)
		return nil, err
	}

	request := artifacts.GetContainerImageRequest{
		ImageId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ArtifactClient.GetContainerImage(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.ContainerImage, nil
}

// Build additional filters
func buildArtifactContainerImageFilters(equalQuals plugin.KeyColumnEqualsQualMap) artifacts.ListContainerImagesRequest {
	request := artifacts.ListContainerImagesRequest{}

	if equalQuals["compartment_id"] != nil {
		request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	if equalQuals["repository_id"] != nil {
		request.RepositoryId = types.String(equalQuals["repository_id"].GetStringValue())
	}

	if equalQuals["repository_name"] != nil {
		request.RepositoryName = types.String(equalQuals["repository_name"].GetStringValue())
	}

	if equalQuals["version"] != nil {
		request.Version = types.String(equalQuals["version"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = types.String(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
