package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerinstances"

	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableContainerInstancesContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_container_instances_container",
		Description:      "Retrieve information about your Containers Running inside Container Instances.",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContainer,
		},
		List: &plugin.ListConfig{
			Hydrate:           listContainers,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
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
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
				{
					Name:    "container_instance_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "isplay name for the Container. Can be renamed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Unique identifier that is immutable on creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "Availability Domain where the Container's Instance is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "container_instance_id",
				Description: "The identifier of the Container Instance on which this container is running.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The time the the Container was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "image_url",
				Description: "The container image information. Currently only support public docker registry. Can be either image name, e.g `containerImage`, image name with version, e.g `containerImage:v1` or complete docker image Url e.g `docker.io/library/containerImage:latest`. If no registry is provided, will default the registry to public docker hub `docker.io/library`. The registry used for container image must be reachable over the Container Instance's VNIC.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "resource_config",
				Description: "The resource config of the Container.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "is_resource_principal_disabled",
				Description: "Determines if the Container will have access to the Container Instance Resource Principal.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "exit_code",
				Description: "The exit code of the container process if it has stopped executing.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getContainer,
			},
			{
				Name:        "command",
				Description: "This command will override the container's entrypoint process. If not specified, the existing entrypoint process defined in the image will be used.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "arguments",
				Description: "A list of string arguments for a Container's entrypoint process.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "additional_capabilities",
				Description: "A list of additional configurable container capabilities.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "volume_mounts",
				Description: "List of the volume mounts.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "health_checks",
				Description: "List of container health checks.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "container_restart_attempt_count",
				Description: "The number of container restart attempts. A restart may be attempted after a health check failure or a container exit, based on the restart policy.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getContainer,
			},
			{
				Name:        "environment_variables",
				Description: "A map of additional environment variables to set in the environment of the container's entrypoint process. These variables are in addition to any variables already defined in the container's image.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainer,
			},
			{
				Name:        "working_directory",
				Description: "The working directory within the Container's filesystem for the Container process. If this is not present, the default working directory from the image will be used.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getContainer,
			},
			{
				Name:        "fault_domain",
				Description: "Fault Domain where the ContainerInstance is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail. For example, can be used to provide actionable information for a resource in Failed state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The time the Container was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "time_terminated",
				Description: "Time at which the container last terminated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeTerminated.Time"),
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
				Transform:   transform.From(containerTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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

func listContainers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listContainers", "Compartment", compartment, "OCI_REGION", region)

	// Return nil, if given compartment_id doesn't match
	if d.EqualsQualString("compartment_id") != "" && compartment != d.EqualsQualString("compartment_id") {
		return nil, nil
	}

	// Create Session
	session, err := containerInstancesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_container_instances_container.listContainers", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request := containerinstances.ListContainersRequest{}
	if d.EqualsQualString("display_name") != "" {
		request.DisplayName = types.String(d.EqualsQualString("display_name"))
	}
	if d.EqualsQualString("lifecycle_state") != "" {
		request.LifecycleState = containerinstances.ContainerLifecycleStateEnum(d.EqualsQualString("lifecycle_state"))
	}
	if d.EqualsQualString("availability_domain") != "" {
		request.AvailabilityDomain = types.String(d.EqualsQualString("availability_domain"))
	}
	if d.EqualsQualString("container_instance_id") != "" {
		request.ContainerInstanceId = types.String(d.EqualsQualString("container_instance_id"))
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
		response, err := session.ContainerInstancesClient.ListContainers(ctx, request)
		if err != nil {
			logger.Error("oci_container_instances_container.listContainers", "api_error", err)
			return nil, err
		}

		for _, ContainerSummary := range response.Items {
			d.StreamListItem(ctx, ContainerSummary)

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

/// HYDRATE FUNCTIONS

func getContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getContainer", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(containerinstances.ContainerSummary).Id
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
	session, err := containerInstancesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_container_instances_container.getContainer", "connection_error", err)
		return nil, err
	}

	request := containerinstances.GetContainerRequest{
		ContainerId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.ContainerInstancesClient.GetContainer(ctx, request)
	if err != nil {
		logger.Error("oci_container_instances_container.getContainer", "api_error", err)
		return nil, err
	}

	return response.Container, nil
}

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func containerTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case containerinstances.ContainerInstanceSummary:
		freeFormTags = d.HydrateItem.(containerinstances.ContainerSummary).FreeformTags
		definedTags = d.HydrateItem.(containerinstances.ContainerSummary).DefinedTags
		systemTags = d.HydrateItem.(containerinstances.ContainerSummary).SystemTags
	case containerinstances.ContainerInstance:
		freeFormTags = d.HydrateItem.(containerinstances.Container).FreeformTags
		definedTags = d.HydrateItem.(containerinstances.Container).DefinedTags
		systemTags = d.HydrateItem.(containerinstances.Container).SystemTags
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
