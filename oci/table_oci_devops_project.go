package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/devops"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableDevopsProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_devops_project",
		Description:      "OCI Devops Project",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getProject,
		},
		List: &plugin.ListConfig{
			Hydrate: listProjects,
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
					Name:    "id",
					Require: plugin.Optional,
				},

			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the project.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "id",
				Description: "The OCID of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the queue was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the queue was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
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

			// // Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(queueTags),
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

func listProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)

	logger.Debug("listProjects", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := devopsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_project.listProjects", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request, isValid := buildProjectFilters(equalQuals, logger)
	if !isValid {
		return nil, nil
	}
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
		response, err := session.DevopsClient.ListProjects(ctx, request)
		// see https://pkg.go.dev/github.com/oracle/oci-go-sdk/v65@v65.28.0/devops#ProjectSummary
		if err != nil {
			logger.Error("oci_devops_project.listProjects", "api_error", err)
			return nil, err
		}

		for _, projectSummary  := range response.Items {
			d.StreamListItem(ctx, projectSummary)

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

//// HYDRATE FUNCTIONS

func getProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getQueue")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getQueue", "Compartment", compartment, "OCI_REGION", region)

	var id string
	// if h.Item != nil {
	// 	id = *h.Item.(queue.QueueSummary).Id
	// } else {
	// 	id = d.KeyColumnQuals["id"].GetStringValue()
	// 	// Restrict the api call to only root compartment/ per region
	// 	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
	// 		return nil, nil
	// 	}
	// }

	if id == "" {
		return nil, nil
	}
	// Create Session
	session, err := devopsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_project.getProject", "connection_error", err)
		return nil, err
	}

	request := devops.GetProjectRequest{
		ProjectId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.DevopsClient.GetProject(ctx, request)
	if err != nil {
		logger.Error("oci_devops_project.getQueue", "api_error", err)
		return nil, err
	}

	return response.Project, nil
}

// Build additional filters
func buildProjectFilters(equalQuals plugin.KeyColumnEqualsQualMap, logger hclog.Logger) (devops.ListProjectsRequest, bool) {
	request := devops.ListProjectsRequest{}
	isValid := true

	if equalQuals["name"] != nil && strings.Trim(equalQuals["name"].GetStringValue(), " ") != "" {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}
	return request, isValid
}

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func projectTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case devops.ProjectSummary:
		freeFormTags = d.HydrateItem.(devops.ProjectSummary).FreeformTags
		definedTags = d.HydrateItem.(devops.ProjectSummary).DefinedTags
	case devops.Project:
		freeFormTags = d.HydrateItem.(devops.Project).FreeformTags
		definedTags = d.HydrateItem.(devops.Project).DefinedTags
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

	return tags, nil
}
