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

func tableContainerInstancesContainerInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_container_instances_container_instance",
		Description: "Retrieve information about your Container Instances.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContainerInstance,
		},
		List: &plugin.ListConfig{
			Hydrate:           listContainerInstances,
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
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "Display name for the Container Instance.",
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
				Description: "Availability Domain where the Container Instance is running.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Container Instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "container_count",
				Description: "The number of containers running on the instance.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "time_created",
				Description: "The time the the Container Instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "shape",
				Description: "The shape of the Container Instance. The shape determines the resources available to the Container Instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shape_config",
				Description: "The shape config of the Container Instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "containers",
				Description: "The Containers on this Instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainerInstance,
			},
			{
				Name:        "vnics",
				Description: "The virtual networks available to containers running on this Container Instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainerInstance,
			},
			{
				Name:        "volumes",
				Description: "A Volume represents a directory with data that is accessible across multiple containers in a ContainerInstance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainerInstance,
			},
			{
				Name:        "dns_config",
				Description: "DNS Config of the container instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainerInstance,
			},
			{
				Name:        "image_pull_secrets",
				Description: "The image pull secrets for accessing private registry to pull images for containers.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getContainerInstance,
			},
			{
				Name:        "container_restart_policy",
				Description: "The container restart policy is applied for all containers in container instance.",
				Type:        proto.ColumnType_STRING,
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
				Description: "The time the ContainerInstance was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "graceful_shutdown_timeout_in_seconds",
				Description: "Duration in seconds processes within a Container have to gracefully terminate. This applies whenever a Container must be halted, such as when the Container Instance is deleted.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "volume_count",
				Description: "The number of volumes that attached to this Instance.",
				Type:        proto.ColumnType_INT,
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
				Transform:   transform.From(containerInstanceTags),
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

func listContainerInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listContainerInstances", "Compartment", compartment, "OCI_REGION", region)

	// Return nil, if given compartment_id doesn't match
	if d.EqualsQualString("compartment_id") != "" && compartment != d.EqualsQualString("compartment_id") {
		return nil, nil
	}

	// Create Session
	session, err := containerInstancesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_container_instances_container_instance.listContainerInstances", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request := containerinstances.ListContainerInstancesRequest{}
	if d.EqualsQualString("display_name") != "" {
		request.DisplayName = types.String(d.EqualsQualString("display_name"))
	}
	if d.EqualsQualString("lifecycle_state") != "" {
		request.LifecycleState = containerinstances.ContainerInstanceLifecycleStateEnum(d.EqualsQualString("lifecycle_state"))
	}
	if d.EqualsQualString("availability_domain") != "" {
		request.AvailabilityDomain = types.String(d.EqualsQualString("availability_domain"))
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
		response, err := session.ContainerInstancesClient.ListContainerInstances(ctx, request)
		if err != nil {
			logger.Error("oci_container_instances_container_instance.listContainerInstances", "api_error", err)
			return nil, err
		}

		for _, containerInstanceSummary := range response.Items {
			d.StreamListItem(ctx, containerInstanceSummary)

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

func getContainerInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getContainerInstance", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(containerinstances.ContainerInstanceSummary).Id
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
		logger.Error("oci_container_instances_container_instance.getContainerInstance", "connection_error", err)
		return nil, err
	}

	request := containerinstances.GetContainerInstanceRequest{
		ContainerInstanceId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.ContainerInstancesClient.GetContainerInstance(ctx, request)
	if err != nil {
		logger.Error("oci_container_instances_container_instance.getContainerInstance", "api_error", err)
		return nil, err
	}

	return response.ContainerInstance, nil
}

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func containerInstanceTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var freeFormTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case containerinstances.ContainerInstanceSummary:
		freeFormTags = d.HydrateItem.(containerinstances.ContainerInstanceSummary).FreeformTags
		definedTags = d.HydrateItem.(containerinstances.ContainerInstanceSummary).DefinedTags
		systemTags = d.HydrateItem.(containerinstances.ContainerInstanceSummary).SystemTags
	case containerinstances.ContainerInstance:
		freeFormTags = d.HydrateItem.(containerinstances.ContainerInstance).FreeformTags
		definedTags = d.HydrateItem.(containerinstances.ContainerInstance).DefinedTags
		systemTags = d.HydrateItem.(containerinstances.ContainerInstance).SystemTags
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
