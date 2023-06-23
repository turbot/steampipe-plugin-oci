package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/devops"

	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableDevopsRepository(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_devops_repository",
		Description:      "OCI Devops Repository",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getRepository,
		},
		List: &plugin.ListConfig{
			Hydrate:           listRepositories,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			// Since we were encountering errors when retrieving results by passing the project_id as an input parameter, we have decided to remove it from the optional key qualifiers.
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
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
				Description: "The OCID of the repository. This value is unique and immutable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_id",
				Description: "The OCID of the DevOps project containing the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Unique name of a repository. This value is mutable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "Tenancy unique namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "project_name",
				Description: "Unique project name in a namespace.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssh_url",
				Description: "SSH URL that you use to git clone, pull and push.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "http_url",
				Description: "HTTP URL that you use to git clone, pull and push.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Details of the repository. Avoid entering confidential information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_branch",
				Description: "The default branch of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "repository_type",
				Description: "Type of repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the repository was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the repository was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the repository.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "branch_count",
				Description: "The count of the branches present in the repository.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getRepository,
			},
			{
				Name:        "commit_count",
				Description: "The count of the commits present in the repository.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getRepository,
			},
			{
				Name:        "size_in_bytes",
				Description: "The size of the repository in bytes.",
				Type:        proto.ColumnType_DOUBLE,
				Hydrate:     getRepository,
			},
			{
				Name:        "mirror_repository_config",
				Description: "Mirror repository configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRepository,
			},
			{
				Name:        "trigger_build_events",
				Description: "Trigger build events supported for this repository.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getRepository,
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
			{
				Name:        "system_tags",
				Description: ColumnDescriptionSystemTags,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(repositoryTags),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listRepositories(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	region := d.EqualsQualString(matrixKeyRegion)

	logger.Debug("listRepositories", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := devOpsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_repository.listRepositories", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request := buildRepositoryFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = oci_common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DevopsClient.ListRepositories(ctx, request)
		if err != nil {
			logger.Error("oci_devops_repository.listRepositories", "api_error", err)
			return nil, err
		}

		for _, repositorySummary := range response.Items {
			d.StreamListItem(ctx, repositorySummary)

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

//// HYDRATE FUNCTIONS

func getRepository(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	region := d.EqualsQualString(matrixKeyRegion)
	logger.Debug("getRepository", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(devops.RepositorySummary).Id
	} else {
		id = d.EqualsQualString("id")
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// check if id is empty
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := devOpsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_repository.getRepository", "connection_error", err)
		return nil, err
	}

	request := devops.GetRepositoryRequest{
		RepositoryId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.DevopsClient.GetRepository(ctx, request)
	if err != nil {
		logger.Error("oci_devops_repository.getQueue", "api_error", err)
		return nil, err
	}

	return response.Repository, nil
}

// Build additional filters
func buildRepositoryFilters(equalQuals plugin.KeyColumnEqualsQualMap) devops.ListRepositoriesRequest {
	request := devops.ListRepositoriesRequest{}

	if equalQuals["name"] != nil && strings.Trim(equalQuals["name"].GetStringValue(), " ") != "" {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil && strings.Trim(equalQuals["lifecycle_state"].GetStringValue(), " ") != "" {
		request.LifecycleState = devops.RepositoryLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}

// Priority order for tags
// 1. Defined Tags
// 2. Free-form Tags
// 3. System Tags
func repositoryTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case devops.RepositorySummary:
		freeFormTags = d.HydrateItem.(devops.RepositorySummary).FreeformTags
		definedTags = d.HydrateItem.(devops.RepositorySummary).DefinedTags
		systemTags = d.HydrateItem.(devops.RepositorySummary).SystemTags
	case devops.Repository:
		freeFormTags = d.HydrateItem.(devops.Repository).FreeformTags
		definedTags = d.HydrateItem.(devops.Repository).DefinedTags
		systemTags = d.HydrateItem.(devops.Repository).SystemTags
	default:
		return nil, nil
	}

	var tags map[string]interface{}
	if freeFormTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeFormTags {
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

	if systemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range systemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}
	return tags, nil
}
