package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_instance",
		Description: "OCI Core Instance",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.AnyColumn([]string{"id"}),
		// 	Hydrate:    getInstance,
		// },
		List: &plugin.ListConfig{
			Hydrate: listCoreInstances,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain the instance is running in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the instance was created",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "dedicated_vm_host_id",
				Description: "The OCID of dedicated VM host.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fault_domain",
				Description: "The name of the fault domain the instance is running in. A fault domain is a grouping of hardware and infrastructure within an availability domain.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ipxe_script",
				Description: "When a bare metal or virtual machine instance boots, the iPXE firmware that runs on the instance is configured to run an iPXE script to continue the boot process.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_mode",
				Description: "Specifies the configuration mode for launching virtual machine (VM) instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_maintenance_reboot_due",
				Description: "The date and time the instance is expected to be stopped/started. After that time if instance hasn't been rebooted, Oracle will reboot the instance within 24 hours of the due time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeMaintenanceRebootDue.time"),
			},
			{
				Name:        "shape",
				Description: "The shape of the instance. The shape determines the number of CPUs and the amount of memory allocated to the instance.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "agent_config",
				Description: "Options for the Oracle Cloud Agent software running on the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "availability_config",
				Description: "Options for defining the availabiity of a VM instance after a maintenance event that impacts the underlying hardware.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "extended_metadata",
				Description: "Additional metadata key/value pairs that user provided to instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "platform_config",
				Description: "The platform configuration for the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata",
				Description: "Custom metadata that you provided to instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "launch_options",
				Description: "LaunchOptions Options for tuning the compatibility and performance of VM shapes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "instance_options",
				Description: "Optional mutable instance options.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shape_config",
				Description: "The shape configuration for an instance. The shape configuration determines the resources allocated to an instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "source_details",
				Description: "Contains the details of the source image for the instance.",
				Type:        proto.ColumnType_JSON,
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
				Description: "Tags added to instances by the service.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(instanceTags),
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
				Transform:   transform.FromCamel().Transform(regionName),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCoreInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("listCoreInstances", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreComputeServiceRegional(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListInstancesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListInstances(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, instance := range response.Items {
			d.StreamListItem(ctx, instance)
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

func getInstance(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getUser")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("oci.getInstance", "Compartment", compartment, "OCI_REGION", region)

	// ➜  steampipe-plugin-oci git:(issue-4) ✗ steampipe query
	// 	Welcome to Steampipe v0.3.0
	// 	For more information, type .help
	// 	> select * from oci_core_instance where id='ocid1.instance.oc1.ap-mumbai-1.anrg6ljr6igdexacgbqlitvrgoq455nhbc3fuim3wpzslnecaqcnfbo4fx4a'
	// 	Error: get call returned 4 results - the key column is not globally unique

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	// session, err := coreComputeService(ctx, d)
	session, err := coreComputeServiceRegional(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetInstanceRequest{
		InstanceId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ComputeClient.GetInstance(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Instance, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func instanceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(core.Instance)

	var tags map[string]interface{}

	if instance.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range instance.FreeformTags {
			tags[k] = v
		}
	}

	if instance.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range instance.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	if instance.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range instance.SystemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

// For the us-phoenix-1 and us-ashburn-1 regions, `phx` and `iad` are returned by ListInstances api, respectively.
// For all other regions, the full region name is returned.
func regionName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	region := types.SafeString(d.Value)

	switch region {
	case "iad":
		return "us-ashburn-1", nil
	case "phx":
		return "us-phoenix-1", nil
	default:
		return region, nil
	}
}
