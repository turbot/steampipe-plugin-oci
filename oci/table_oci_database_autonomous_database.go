package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOciDatabaseAutonomousDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_autonomous_database",
		Description: "OCI Database Autonomous Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAutonomousDatabase,
		},
		List: &plugin.ListConfig{
			Hydrate: listAutonomousDatabases,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "autonomous_container_database_id",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "db_version",
					Require: plugin.Optional,
				},
				{
					Name:    "db_workload",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "infrastructure_type",
					Require: plugin.Optional,
				},
				{
					Name:      "is_data_guard_enabled",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:      "is_free_tier",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The user-friendly name for the Autonomous Database. The name does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_name",
				Description: "The database name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the Autonomous Database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the Autonomous Database was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "are_primary_whitelisted_ips_used",
				Description: "This field will be null if the Autonomous Database is not Data Guard enabled or Access Control is disabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "autonomous_container_database_id",
				Description: "The Autonomous Container Database OCID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "cpu_core_count",
				Description: "The number of OCPU cores to be made available to the database.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuCoreCount"),
			},
			{
				Name:        "data_safe_status",
				Description: "Status of the Data Safe registration for this Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_storage_size_in_gbs",
				Description: "The quantity of data in the database, in gigabytes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DataStorageSizeInGBs"),
			},
			{
				Name:        "data_storage_size_in_tbs",
				Description: "The quantity of data in the database, in terabytes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DataStorageSizeInTBs"),
			},
			{
				Name:        "db_version",
				Description: "A valid Oracle Database version for Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_workload",
				Description: "The Autonomous Database workload type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "failed_data_recovery_in_seconds",
				Description: "Indicates the number of seconds of data loss for a Data Guard failover.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "infrastructure_type",
				Description: "The infrastructure type this resource belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_access_control_enabled",
				Description: "Indicates if the database-level access control is enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_auto_scaling_enabled",
				Description: "Indicates if auto scaling is enabled for the Autonomous Database CPU core count.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_data_guard_enabled",
				Description: "Indicates whether the Autonomous Database has Data Guard enabled.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_dedicated",
				Description: "True if the database uses dedicated Exadata infrastructure.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_free_tier",
				Description: "Indicates if this is an Always Free resource. The default value is false.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_preview",
				Description: "Indicates if the Autonomous Database version is a preview version.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_refreshable_clone",
				Description: "Indicates whether the Autonomous Database is a refreshable clone.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "key_store_id",
				Description: "The OCID of the key store.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "key_store_wallet_name",
				Description: "The wallet name for Oracle Key Vault.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the key container that is used as the master encryption key in database transparent data encryption (TDE) operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KmsKeyId"),
			},
			{
				Name:        "kms_key_version_id",
				Description: "The OCID of the key container version that is used in database transparent data encryption (TDE) operations KMS Key can have multiple key versions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KmsKeyVersionId"),
			},
			{
				Name:        "lifecycle_details",
				Description: "Information about the current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "The Oracle license model that applies to the Oracle Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "open_mode",
				Description: "The `DATABASE OPEN` mode. You can open the database in `READ_ONLY` or `READ_WRITE` mode.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operations_insights_status",
				Description: "Status of Operations Insights for this Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permission_level",
				Description: "The Autonomous Database permission level.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_endpoint",
				Description: "The private endpoint for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_endpoint_ip",
				Description: "The private endpoint Ip address for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_endpoint_label",
				Description: "The private endpoint label for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "refreshable_mode",
				Description: "The refresh mode of the clone. AUTOMATIC indicates that the clone is automatically being refreshed with data from the source Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "refreshable_status",
				Description: "The refresh status of the clone. REFRESHING indicates that the clone is currently being refreshed with data from the source Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The role of the Autonomous Data Guard-enabled Autonomous Container Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_console_url",
				Description: "The URL of the Service Console for the Autonomous Database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_id",
				Description: "The OCID of the source Autonomous Database that was cloned to create the current Autonomous Database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet the resource is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_deletion_of_free_autonomous_database",
				Description: "The date and time the Always Free database will be automatically deleted because of inactivity. If the database is in the STOPPED state and without activity until this time, it will be deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeDeletionOfFreeAutonomousDatabase.Time"),
			},
			{
				Name:        "time_maintenance_begin",
				Description: "The date and time when maintenance will begin.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeMaintenanceBegin.Time"),
			},
			{
				Name:        "time_maintenance_end",
				Description: "The date and time when maintenance will end.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeMaintenanceEnd.Time"),
			},
			{
				Name:        "time_of_last_failover",
				Description: "The timestamp of the last failover operation.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfLastFailover.Time"),
			},
			{
				Name:        "time_of_last_refresh",
				Description: "The date and time when last refresh happened.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfLastRefresh.Time"),
			},
			{
				Name:        "time_of_last_refresh_point",
				Description: "The refresh point timestamp (UTC). The refresh point is the time to which the database was most recently refreshed. Data created after the refresh point is not included in the refresh.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfLastRefreshPoint.Time"),
			},
			{
				Name:        "time_of_last_switchover",
				Description: "The timestamp of the last switchover operation for the Autonomous Database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfLastSwitchover.Time"),
			},
			{
				Name:        "time_of_next_refresh",
				Description: "The date and time of next refresh.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfNextRefresh.Time"),
			},
			{
				Name:        "time_reclamation_of_free_autonomous_database",
				Description: "The date and time the Always Free database will be stopped because of inactivity. If this time is reached without any database activity, the database will automatically be put into the STOPPED state.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeReclamationOfFreeAutonomousDatabase.Time"),
			},
			{
				Name:        "used_data_storage_size_in_tbs",
				Description: "The amount of storage that has been used, in terabytes.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("UsedDataStorageSizeInTBs"),
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the Oracle Cloud Infrastructure vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VaultId"),
			},
			{
				Name:        "apex_details",
				Description: "Information about Oracle APEX Application Development.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "available_upgrade_versions",
				Description: "List of Oracle Database versions available for a database upgrade. If there are no version upgrades available, this list is empty.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "backup_config",
				Description: "Autonomous Database configuration details for storing manual backups in the Object Storage service.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connection_strings",
				Description: "The connection string used to connect to the Autonomous Database. The username for the Service Console is ADMIN. Use the password you entered when creating the Autonomous Database for the password value.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "connection_urls",
				Description: "The URLs for accessing Oracle Application Express (APEX) and SQL Developer Web with a browser from a Compute instance within your VCN or that has a direct connection to your VCN. Note that these URLs are provided by the console only for databases on dedicated Exadata infrastructure.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nsg_ids",
				Description: "A list of the OCIDs of the network security groups (NSGs) that this resource belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "standby_db",
				Description: "Autonomous Data Guard standby database details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "standby_whitelisted_ips",
				Description: "The client IP access control list (ACL). This feature is available for autonomous databases on shared Exadata infrastructure and on Exadata Cloud@Customer. Only clients connecting from an IP address included in the ACL may access the Autonomous Database instance. For shared Exadata infrastructure, this is an array of CIDR (Classless Inter-Domain Routing) notations for a subnet or VCN OCID.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "whitelisted_ips",
				Description: "The client IP access control list (ACL).",
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
				Description: ColumnDescriptionSystemTags,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(autonomousDatabaseTags),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAutonomousDatabases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listAutonomousDatabases", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	quals := d.Quals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildAutonomousDatabaseFilter(equalQuals, quals)
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
		response, err := session.DatabaseClient.ListAutonomousDatabases(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, database := range response.Items {
			d.StreamListItem(ctx, database)

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

func getAutonomousDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getAutonomousDatabase", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := database.GetAutonomousDatabaseRequest{
		AutonomousDatabaseId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetAutonomousDatabase(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.AutonomousDatabase, nil
}

//// TRANSFORM FUNCTION

func autonomousDatabaseTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.AutonomousDatabaseSummary:
		database := d.HydrateItem.(database.AutonomousDatabaseSummary)
		freeformTags = database.FreeformTags
		definedTags = database.DefinedTags
		systemTags = database.SystemTags
	case database.AutonomousDatabase:
		database := d.HydrateItem.(database.AutonomousDatabase)
		freeformTags = database.FreeformTags
		definedTags = database.DefinedTags
		systemTags = database.SystemTags
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

func buildAutonomousDatabaseFilter(equalQuals plugin.KeyColumnEqualsQualMap, quals plugin.KeyColumnQualMap) database.ListAutonomousDatabasesRequest {
	request := database.ListAutonomousDatabasesRequest{}

	filterQuals := []string{
		"autonomous_container_database_id",
		"db_version",
		"db_workload",
		"display_name",
		"infrastructure_type",
		"is_data_guard_enabled",
		"is_free_tier",
		"lifecycle_state",
	}

	for _, columnName := range filterQuals {
		if equalQuals[columnName] != nil {
			switch columnName {
			case "autonomous_container_database_id":
				request.AutonomousContainerDatabaseId = types.String(equalQuals[columnName].GetStringValue())
			case "db_version":
				request.DbVersion = types.String(equalQuals[columnName].GetStringValue())
			case "db_workload":
				request.DbWorkload = database.AutonomousDatabaseSummaryDbWorkloadEnum(equalQuals[columnName].GetStringValue())
			case "display_name":
				request.DisplayName = types.String(equalQuals[columnName].GetStringValue())
			case "infrastructure_type":
				request.InfrastructureType = database.AutonomousDatabaseSummaryInfrastructureTypeEnum(equalQuals[columnName].GetStringValue())
			case "lifecycle_state":
				request.LifecycleState = database.AutonomousDatabaseSummaryLifecycleStateEnum(equalQuals[columnName].GetStringValue())
			case "is_data_guard_enabled":
				request.IsDataGuardEnabled = types.Bool(equalQuals[columnName].GetBoolValue())
			case "is_free_tier":
				request.IsFreeTier = types.Bool(equalQuals[columnName].GetBoolValue())
			}
		}
	}

	boolNEQuals := []string{
		"is_data_guard_enabled",
		"is_free_tier",
	}
	// Non-Equals Qual Map handling
	for _, qual := range boolNEQuals {
		if quals[qual] != nil {
			for _, q := range quals[qual].Quals {
				value := q.Value.GetBoolValue()
				if q.Operator == "<>" {
					if qual == "is_data_guard_enabled" {
						request.IsDataGuardEnabled = types.Bool(false)
						if !value {
							request.IsDataGuardEnabled = types.Bool(true)
						}
					}
					if qual == "is_free_tier" {
						request.IsFreeTier = types.Bool(false)
						if !value {
							request.IsFreeTier = types.Bool(true)
						}
					}
				}
			}
		}
	}
	return request
}
