package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableMySQLDBSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_db_system",
		Description: "OCI MySQL DB System",
		List: &plugin.ListConfig{
			Hydrate: listMySQLDBSystems,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "configuration_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "is_analytics_cluster_attached",
					Require: plugin.Optional,
				},
				{
					Name:    "is_heat_wave_cluster_attached",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getMySQLDBSystem,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The user-friendly name for the DB System. It does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the DB System.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the DB System.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "configuration_id",
				Description: "The OCID of the Configuration to be used for Instances in this DB System.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet the DB System is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "time_created",
				Description: "The date and time the DB System was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "availability_domain",
				Description: "The Availability Domain where the primary DB System should be located.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_storage_size_in_gbs",
				Description: "Initial size of the data volume in GiBs that will be created and attached.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getMySQLDBSystem,
				Transform:   transform.FromField("DataStorageSizeInGBs"),
			},
			{
				Name:        "description",
				Description: "User-provided data about the DB System.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fault_domain",
				Description: "The name of the fault domain the DB System is located in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "hostname_label",
				Description: "The hostname for the primary endpoint of the DB System.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "ip_address",
				Description: "The IP address the DB System is configured to listen on.",
				Type:        proto.ColumnType_IPADDR,
				Hydrate:     getMySQLDBSystem,
				Transform:   transform.FromField("IpAddress"),
			},
			{
				Name:        "is_analytics_cluster_attached",
				Description: "If the DB System has an Analytics Cluster attached.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_heat_wave_cluster_attached",
				Description: "Whether the DB System has a HeatWave cluster attached.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycleState.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "mysql_version",
				Description: "Name of the MySQL Version in use for the DB System.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "port",
				Description: "The port for primary endpoint of the DB System to listen on.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "port_x",
				Description: "The network port on which X Plugin listens for TCP/IP connections. This is the X Plugin equivalent of port.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "shape_name",
				Description: "The shape of the primary instances of the DB System.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "time_updated",
				Description: "The time the DB System was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			// json fields
			{
				Name:        "analytics_cluster",
				Description: "A summary of an Analytics Cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "backup_policy",
				Description: "BackupPolicy The Backup policy for the DB System.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "channels",
				Description: "A list with a summary of all the Channels attached to the DB System.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "endpoints",
				Description: "The network endpoints available for this DB System.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "maintenance",
				Description: "The Maintenance Policy for the DB System.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLDBSystem,
			},
			{
				Name:        "source",
				Description: "DbSystemSource Parameters detailing how to provision the initial data of the DB System.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLDBSystem,
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
				Transform:   transform.From(dbSystemTags),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listMySQLDBSystems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listMySQLDBSystems", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := mySQLDBSystemService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildMySQLDBSystemFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
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
		response, err := session.MySQLDBSystemClient.ListDbSystems(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, dbSystem := range response.Items {
			d.StreamListItem(ctx, dbSystem)

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

func getMySQLDBSystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getMySQLDBSystem", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(mysql.DbSystemSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := mySQLDBSystemService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.GetDbSystemRequest{
		DbSystemId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.MySQLDBSystemClient.GetDbSystem(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DbSystem, nil
}

//// TRANSFORM FUNCTION

func dbSystemTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := dbSystemFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := dbSystemDefinedTags(d.HydrateItem)

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

func dbSystemFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case mysql.DbSystem:
		return item.FreeformTags
	case mysql.DbSystemSummary:
		return item.FreeformTags
	}
	return nil
}

func dbSystemDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case mysql.DbSystem:
		return item.DefinedTags
	case mysql.DbSystemSummary:
		return item.DefinedTags
	}
	return nil
}

// Build additional filters
func buildMySQLDBSystemFilters(equalQuals plugin.KeyColumnEqualsQualMap) mysql.ListDbSystemsRequest {
	request := mysql.ListDbSystemsRequest{}

	if equalQuals["configuration_id"] != nil {
		request.ConfigurationId = types.String(equalQuals["configuration_id"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["id"] != nil {
		request.DbSystemId = types.String(equalQuals["id"].GetStringValue())
	}
	if equalQuals["is_analytics_cluster_attached"] != nil {
		request.IsAnalyticsClusterAttached = types.Bool(equalQuals["is_analytics_cluster_attached"].GetBoolValue())
	}
	if equalQuals["is_heat_wave_cluster_attached"] != nil {
		request.IsHeatWaveClusterAttached = types.Bool(equalQuals["is_heat_wave_cluster_attached"].GetBoolValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = mysql.DbSystemLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
