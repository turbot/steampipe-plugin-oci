package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/devops"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableDevopsProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_devops_project",
		Description:      "OCI Devops Project",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDevopsProject,
		},
		List: &plugin.ListConfig{
			Hydrate: listDevopsProjects,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
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
			},
			{
				Name:        "description",
				Description: "The description of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the project was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the project was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail. For example, can be used to provide actionable information for a resource in Failed state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "Namespace associated with the project.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "notification_topic_id",
				Description: "The topic ID for the topic where project notifications will be published to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(transformNotificationTopic),
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
				Transform:   transform.From(devopsProjectTags),
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

func listDevopsProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	region := d.EqualsQualString(matrixKeyRegion)

	logger.Debug("listDevopsProjects", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := devopsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_project.listDevopsProjects", "connection_error", err)
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
			logger.Error("oci_devops_project.listDevopsProjects", "api_error", err)
			return nil, err
		}

		for _, projectSummary := range response.Items {
			d.StreamListItem(ctx, projectSummary)

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

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getDevopsProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	id := d.EqualsQualString("id")

	// Restrict the api call to only root compartment
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := devopsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_devops_project.getDevopsProject", "connection_error", err)
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
		logger.Error("oci_devops_project.getDevopsProject", "api_error", err)
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

func transformNotificationTopic(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var topicId string
	switch d.HydrateItem.(type) {
	case devops.ProjectSummary:
		topicId = *d.HydrateItem.(devops.ProjectSummary).NotificationConfig.TopicId
	case devops.Project:
		topicId = *d.HydrateItem.(devops.Project).NotificationConfig.TopicId
	}
	return topicId, nil
}

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func devopsProjectTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case devops.ProjectSummary:
		freeFormTags = d.HydrateItem.(devops.ProjectSummary).FreeformTags
		definedTags = d.HydrateItem.(devops.ProjectSummary).DefinedTags
		systemTags = d.HydrateItem.(devops.ProjectSummary).SystemTags
	case devops.Project:
		freeFormTags = d.HydrateItem.(devops.Project).FreeformTags
		definedTags = d.HydrateItem.(devops.Project).DefinedTags
		systemTags = d.HydrateItem.(devops.Project).SystemTags
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
