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

func tableArtifactGenericArtifact(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifact_generic_artifact",
		Description:      "OCI Generic Artifact",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getArtifactGenericArtifact,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listArtifactRepositories,
			Hydrate:       listArtifactGenericArtifacts,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "repository_id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "artifact_path",
					Require: plugin.Optional,
				},
				{
					Name:    "version",
					Require: plugin.Optional,
				},
				{
					Name:    "sha256",
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
			// The column 'display_name' has been renamed to 'name' in this table due to its relationship with the table 'oci_artifact_repository'.
			// The table 'oci_artifact_repository' has an optional qualifier called 'display_name' in its list configuration.
			// Additionally, this table has a property named 'DisplayName' which can also be used as an optional key qualifier.
			// When running a query like 'select * from oci_artifact_generic_artifact where display_name = 'test/artifact:1'', the value of 'display_name' is passed as a parameter to the 'listArtifactRepositories' API of the 'oci_artifact_repository' table, which is its parent table.
			// However, this results in an error when the 'listArtifactRepositories' function is called, with the following details: Error: Error returned by Artifacts Service. Http Status Code: 400. Error Code: BadRequest. Opc request id: c1f7fce0f00b215e711b20e9c28c58c6/2bfb17a9e4750c67c6b384c1. Message: Repository name invalid: 'test/artifact:1'.
			// The operation name associated with this error is 'ListRepositories'.
			{
				Name:        "name",
				Description: "The artifact name with the format of `<artifact-path>:<artifact-version>`. The artifact name is truncated to a maximum length of 255.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the artifact.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Generic Artifact was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "repository_id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "artifact_path",
				Description: "A user-defined path to describe the location of an artifact. Slashes do not create a directory structure, but you can use slashes to organize the repository. An artifact path does not include an artifact version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "A user-defined string to describe the artifact version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sha256",
				Description: "The SHA256 digest for the artifact. When you upload an artifact to the repository, a SHA256 digest is calculated and added to the artifact properties.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size_in_bytes",
				Description: "The size of the artifact in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the artifact.",
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
				Transform:   transform.From(artifactGenericArtifactTags),
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

//// LIST FUNCTION

func listArtifactGenericArtifacts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	repository := h.Item.(artifacts.RepositorySummary)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifact_generic_artifact.listArtifactGenericArtifacts", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Return nil, if given repository_id doesn't match
	if equalQuals["repository_id"] != nil && *repository.GetId() != equalQuals["repository_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifact_generic_artifact.listArtifactGenericArtifacts", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildArtifactGenericArtifactFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.RepositoryId = repository.GetId()
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
		response, err := session.ArtifactClient.ListGenericArtifacts(ctx, request)
		if err != nil {
			logger.Error("oci_artifact_generic_artifact.listArtifactGenericArtifacts", "api_error", err)
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

func getArtifactGenericArtifact(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifact_generic_artifact.getArtifactGenericArtifact", "Compartment", compartment, "OCI_REGION", region)

	id := d.EqualsQuals["id"].GetStringValue()
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifact_generic_artifact.getArtifactGenericArtifact", "connection_error", err)
		return nil, err
	}

	request := artifacts.GetGenericArtifactRequest{
		ArtifactId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ArtifactClient.GetGenericArtifact(ctx, request)
	if err != nil {
		logger.Error("oci_artifact_generic_artifact.getArtifactGenericArtifact", "api_error", err)
		return nil, err
	}
	return response.GenericArtifact, nil
}

//// TRANSFORM FUNCTION

func artifactGenericArtifactTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case artifacts.GenericArtifact:
		obj := d.HydrateItem.(artifacts.GenericArtifact)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case artifacts.GenericArtifactSummary:
		obj := d.HydrateItem.(artifacts.GenericArtifactSummary)
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
func buildArtifactGenericArtifactFilters(equalQuals plugin.KeyColumnEqualsQualMap) artifacts.ListGenericArtifactsRequest {
	request := artifacts.ListGenericArtifactsRequest{}

	if equalQuals["name"] != nil {
		request.DisplayName = types.String(equalQuals["name"].GetStringValue())
	}

	if equalQuals["artifact_path"] != nil {
		request.ArtifactPath = types.String(equalQuals["artifact_path"].GetStringValue())
	}

	if equalQuals["version"] != nil {
		request.Version = types.String(equalQuals["version"].GetStringValue())
	}

	if equalQuals["sha256"] != nil {
		request.Sha256 = types.String(equalQuals["sha256"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = types.String(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
