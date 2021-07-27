package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableMySQLConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_configuration",
		Description: "OCI My SQL Configuration",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getConfiguration,
		},
		List: &plugin.ListConfig{
			Hydrate: listMySQLConfigurations,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The display name of the Configuration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the Configuration.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "description",
				Description: "The current state of the Configuration",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Configuration",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that budget was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "parent_configuration_id",
				Description: "The OCID of the Configuration from which this Configuration is derived.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getConfiguration,
				Transform:   transform.FromField("ParentConfigurationId"),
			},
			{
				Name:        "shape_name",
				Description: "The name of the associated Shape.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The date and time the Configuration was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "type",
				Description: "The Configuration type, DEFAULT or CUSTOM..",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "variables",
				Description: "ConfigurationVariables User controllable service variables.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getConfiguration,
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

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(configurationTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// OCI standard columns
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listMySQLConfigurations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listMySQLConfigurations", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := mySQLaasService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.ListConfigurationsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.MySQLaasClient.ListConfigurations(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, configuration := range response.Items {
			d.StreamListItem(ctx, configuration)
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

func getConfiguration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getConfiguration", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		configuration := h.Item.(mysql.ConfigurationSummary)
		id = *configuration.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// handle empty id in get call
		if strings.TrimSpace(id) == "" {
			return nil, nil
		}
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// Create Session
	session, err := mySQLaasService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.GetConfigurationRequest{
		ConfigurationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.MySQLaasClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Configuration, nil
}

//// TRANSFORM FUNCTION

func configurationTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case mysql.Configuration:
		configuration := d.HydrateItem.(mysql.Configuration)
		freeformTags = configuration.FreeformTags
		definedTags = configuration.DefinedTags
	case mysql.ConfigurationSummary:
		configuration := d.HydrateItem.(mysql.ConfigurationSummary)
		freeformTags = configuration.FreeformTags
		definedTags = configuration.DefinedTags
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
