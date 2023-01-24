package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreInstanceConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_instance_configuration",
		Description: "OCI Core Instance Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getInstanceConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreInstanceConfigurations,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the instance configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the instance configuration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "deferred_fields",
				Description: "Parameters that were not specified when the instance configuration was created, but that are required to launch an instance from the instance configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInstanceConfiguration,
			},
			{
				Name:        "instance_details",
				Description: "The intance configuration details.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getInstanceConfiguration,
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(instanceConfigurationTags),
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
				Transform:   transform.FromCamel().Transform(extractInstanceConfigurationRegion),
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

func listCoreInstanceConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	// Create Session
	session, err := coreComputeManagementService(ctx, d, region)
	if err != nil {
		logger.Error("oci_core_instance_configuration.listCoreInstanceConfigurations", "connection_error", err)
		return nil, err
	}

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	
	request := core.ListInstanceConfigurationsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
	}
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
		response, err := session.ComputeManagementClient.ListInstanceConfigurations(ctx, request)
		if err != nil {
			logger.Error("oci_core_instance_configuration.listCoreInstanceConfigurations", "api_error", err)
			return nil, err
		}

		for _, configuration := range response.Items {
			d.StreamListItem(ctx, configuration)

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

//// HYDRATE FUNCTION

func getInstanceConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// Return nil, if no input provided
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeManagementService(ctx, d, region)
	if err != nil {
		logger.Error("oci_core_instance_configuration.getInstanceConfiguration", "connection_error", err)
		return nil, err
	}

	request := core.GetInstanceConfigurationRequest{
		InstanceConfigurationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ComputeManagementClient.GetInstanceConfiguration(ctx, request)
	if err != nil {
		logger.Error("oci_core_instance_configuration.getInstanceConfiguration", "api_error", err)
		return nil, err
	}

	return response.InstanceConfiguration, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func instanceConfigurationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var tags map[string]interface{}

	switch d.HydrateItem.(type) {
	case core.InstanceConfiguration:
		configuration := d.HydrateItem.(core.InstanceConfiguration)
		if configuration.FreeformTags != nil {
			tags = map[string]interface{}{}
			for k, v := range configuration.FreeformTags {
				tags[k] = v
			}
		}

		if configuration.DefinedTags != nil {
			if tags == nil {
				tags = map[string]interface{}{}
			}
			for _, v := range configuration.DefinedTags {
				for key, value := range v {
					tags[key] = value
				}

			}
		}
	case core.InstanceConfigurationSummary:
		configuration := d.HydrateItem.(core.InstanceConfigurationSummary)
		if configuration.FreeformTags != nil {
			tags = map[string]interface{}{}
			for k, v := range configuration.FreeformTags {
				tags[k] = v
			}
		}

		if configuration.DefinedTags != nil {
			if tags == nil {
				tags = map[string]interface{}{}
			}
			for _, v := range configuration.DefinedTags {
				for key, value := range v {
					tags[key] = value
				}

			}
		}

	}

	return tags, nil
}

// For the us-phoenix-1 and us-ashburn-1 regions, `phx` and `iad` are returned by ListInstances api, respectively.
// For all other regions, the full region name is returned.
func extractInstanceConfigurationRegion(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var id string
	switch d.HydrateItem.(type) {
	case core.InstanceConfiguration:
		id = *d.HydrateItem.(core.InstanceConfiguration).Id
	case core.InstanceConfigurationSummary:
		id = *d.HydrateItem.(core.InstanceConfigurationSummary).Id
	}

	return strings.Split(id, ".")[3], nil
}
