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

func tableOciDatabaseDBSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_db_system",
		Description: "OCI Database DB System",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabaseDBSystem,
		},
		List: &plugin.ListConfig{
			Hydrate: listDatabaseDBSystems,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
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
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The user-friendly name for the DB system. The name does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The name of the availability domain that the DB system is located in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the DB system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the DB system was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "backup_subnet_id",
				Description: "The OCID of the backup network subnet the DB system is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "cluster_name",
				Description: "The cluster name for exadata and 2-node RAC virtual machine DB systems.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cpu_core_count",
				Description: "The number of CPU cores enabled on the DB system.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuCoreCount"),
			},
			{
				Name:        "database_edition",
				Description: "The oracle database edition that applies to all the databases on the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_storage_percentage",
				Description: "The percentage assigned to data storage.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "data_storage_size_in_gbs",
				Description: "The data storage size, in gigabytes, that is currently available to the DB system.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DataStorageSizeInGBs"),
			},
			{
				Name:        "db_system_options_storage_management",
				Description: "The storage option used in DB system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DbSystemOptions.StorageManagement"),
			},
			{
				Name:        "disk_redundancy",
				Description: "The type of redundancy configured for the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain",
				Description: "The domain name for the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "host_name",
				Description: "The hostname for the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the key container that is used as the master encryption key in database transparent data encryption (TDE) operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "last_maintenance_run_id",
				Description: "The OCID of the last maintenance run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "last_patch_history_entry_id",
				Description: "The OCID of the last patch history.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "license_model",
				Description: "The oracle license model that applies to all the databases on the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "listener_port",
				Description: "The port number configured for the listener on the DB system.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "next_maintenance_run_id",
				Description: "The OCID of the next maintenance run.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "node_count",
				Description: "The number of nodes in the DB system.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "point_in_time_data_disk_clone_timestamp",
				Description: "The point in time for a cloned database system when the data disks were cloned from the source database system.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("PointInTimeDataDiskCloneTimestamp.Time"),
			},
			{
				Name:        "reco_storage_size_in_gb",
				Description: "The RECO/REDO storage size, in gigabytes, that is currently allocated to the DB system.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("RecoStorageSizeInGB"),
			},
			{
				Name:        "scan_dns_name",
				Description: "The FQDN of the DNS record for the SCAN IP addresses that are associated with the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scan_dns_record_id",
				Description: "The OCID of the DNS record for the SCAN IP addresses that are associated with the DB system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "source_db_system_id",
				Description: "The OCID of the DB system from where the DB system is created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "sparse_diskgroup",
				Description: "True, If sparse diskgroup is configured for exadata DB system.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet the DB system is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_zone",
				Description: "The time zone of the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version",
				Description: "The oracle database version of the DB system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone_id",
				Description: "The OCID of the zone the DB system is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "backup_network_nsg_ids",
				Description: "A list of the OCIDs of the network security groups (NSGs) that the backup network of this DB system belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "fault_domains",
				Description: "List of the fault domains in which this DB system is provisioned.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iorm_config_cache",
				Description: "The IORM configuration of the DB system.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDatabaseDBSystem,
			},
			{
				Name:        "maintenance_window",
				Description: "The maintenance window of the DB system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "nsg_ids",
				Description: "A list of the OCIDs of the network security groups (NSGs) that this resource belongs to.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "scan_ip_ids",
				Description: "The OCID of the single client access name (SCAN) IP addresses associated with the DB system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssh_public_keys",
				Description: "The public key portion of one or more key pairs used for SSH access to the DB system.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vip_ids",
				Description: "A list of the OCIDs of the virtual IP (VIP) addresses associated with the DB system.",
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

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(databaseDBSystemTags),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listDatabaseDBSystems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listDatabaseDBSystems", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

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
	request := buildDatabaseDBSystemFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
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
		response, err := session.DatabaseClient.ListDbSystems(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, dbSystem := range response.Items {
			d.StreamListItem(ctx, dbSystem)

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

func getDatabaseDBSystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getDatabaseDBSystem", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(database.DbSystemSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
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
	session, err := databaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := database.GetDbSystemRequest{
		DbSystemId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetDbSystem(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DbSystem, nil
}

//// TRANSFORM FUNCTION

func databaseDBSystemTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.DbSystemSummary:
		database := d.HydrateItem.(database.DbSystemSummary)
		freeformTags = database.FreeformTags
		definedTags = database.DefinedTags
	case database.DbSystem:
		database := d.HydrateItem.(database.DbSystem)
		freeformTags = database.FreeformTags
		definedTags = database.DefinedTags
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

// Build additional filters
func buildDatabaseDBSystemFilters(equalQuals plugin.KeyColumnEqualsQualMap) database.ListDbSystemsRequest {
	request := database.ListDbSystemsRequest{}

	if equalQuals["availability_domain"] != nil {
		request.AvailabilityDomain = types.String(equalQuals["availability_domain"].GetStringValue())
	}
	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = database.DbSystemSummaryLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
