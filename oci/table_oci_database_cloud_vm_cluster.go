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

func tableOciDatabaseCloudVMCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_cloud_vm_cluster",
		Description: "OCI Database Cloud VM Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabaseCloudVmCluster,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listDatabaseCloudVMClusters,
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
					Name:    "cloud_exadata_infrastructure_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "cluster_name",
				Description: "The cluster name for cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ClusterName"),
			},
			{
				Name:        "display_name",
				Description: "The user-friendly name for the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LifecycleState"),
			},
			{
				Name:        "id",
				Description: "The OCID of the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "availability_domain",
				Description: "The name of the availability domain that the cloud Exadata infrastructure resource is located in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AvailabilityDomain"),
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet associated with the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubnetId"),
			},
			{
				Name:        "shape",
				Description: "The model name of the Exadata hardware running the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Shape"),
			},
			{
				Name:        "hostname",
				Description: "The hostname for the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Hostname"),
			},
			{
				Name:        "domain",
				Description: "The domain name for the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Domain"),
			},
			{
				Name:        "cpu_core_count",
				Description: "The number of CPU cores enabled on the cloud VM cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpuCoreCount"),
			},
			{
				Name:        "cloud_exadata_infrastructure_id",
				Description: "The OCID of the cloud Exadata infrastructure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CloudExadataInfrastructureId"),
			},
			{
				Name:        "backup_subnet_id",
				Description: "The OCID of the backup network subnet associated with the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackupSubnetId"),
			},
			{
				Name:        "last_update_history_entry_id",
				Description: "The OCID of the last maintenance update history entry. This value is updated when a maintenance update starts.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastUpdateHistoryEntryId"),
			},
			{
				Name:        "listener_port",
				Description: "The port number configured for the listener on the cloud VM cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ListenerPort"),
			},
			{
				Name:        "node_count",
				Description: "The number of nodes in the cloud VM cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("NodeCount"),
			},
			{
				Name:        "storage_size_in_gbs",
				Description: "The storage allocation for the disk group, in gigabytes (GB).",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("StorageSizeInGBs"),
			},
			{
				Name:        "time_created",
				Description: "The date and time that the cloud VM cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycle state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LifecycleDetails"),
			},
			{
				Name:        "time_zone",
				Description: "The time zone of the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeZone"),
			},
			{
				Name:        "ocpu_count",
				Description: "The number of OCPU cores to enable on the cloud VM cluster.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("OcpuCount"),
			},
			{
				Name:        "memory_size_in_gbs",
				Description: "The memory to be allocated in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MemorySizeInGBs"),
			},
			{
				Name:        "db_node_storage_size_in_gbs",
				Description: "The local node storage to be allocated in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DbNodeStorageSizeInGBs"),
			},
			{
				Name:        "data_storage_size_in_tbs",
				Description: "The data disk group size to be allocated in TBs.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DataStorageSizeInTBs"),
			},
			{
				Name:        "data_storage_percentage",
				Description: "The percentage assigned to DATA storage.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DataStoragePercentage"),
			},
			{
				Name:        "is_local_backup_enabled",
				Description: "If true, database backup on local Exadata storage is configured for the cloud VM cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsLocalBackupEnabled"),
			},
			{
				Name:        "is_sparse_diskgroup_enabled",
				Description: "If true, sparse disk group is configured for the cloud VM cluster. If false, sparse disk group is not created.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsSparseDiskgroupEnabled"),
			},
			{
				Name:        "gi_version",
				Description: "A valid Oracle Grid Infrastructure (GI) software version.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("GiVersion"),
			},
			{
				Name:        "system_version",
				Description: "Operating system version of the image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SystemVersion"),
			},
			{
				Name:        "license_model",
				Description: "The Oracle license model that applies to the cloud VM cluster. The default is LICENSE_INCLUDED.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LicenseModel"),
			},
			{
				Name:        "disk_redundancy",
				Description: "The type of redundancy configured for the cloud Vm cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DiskRedundancy"),
			},
			{
				Name:        "scan_dns_record_id",
				Description: "The OCID of the DNS record for the SCAN IP addresses associated with the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ScanDnsRecordId"),
			},
			{
				Name:        "scan_dns_name",
				Description: "The FQDN of the DNS record for the SCAN IP addresses associated with the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ScanDnsName"),
			},
			{
				Name:        "zone_id",
				Description: "The OCID of the zone the cloud VM cluster is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneId"),
			},
			{
				Name:        "scan_listener_port_tcp",
				Description: "The TCP Single Client Access Name (SCAN) port. The default port is 1521.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ScanListenerPortTcp"),
			},
			{
				Name:        "scan_listener_port_tcp_ssl",
				Description: "The TCPS Single Client Access Name (SCAN) port. The default port is 2484.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ScanListenerPortTcpSsl"),
			},

			// JSON fields
			{
				Name:        "ssh_public_keys",
				Description: "The public key portion of one or more key pairs used for SSH access to the cloud VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SshPublicKeys"),
			},
			{
				Name:        "nsg_ids",
				Description: "The list of OCIDs for the network security groups (NSGs) to which this resource belongs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NsgIds"),
			},
			{
				Name:        "backup_network_nsg_ids",
				Description: "A list of the OCIDs of the network security groups (NSGs) that the backup network of this DB system belongs to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackupNetworkNsgIds"),
			},
			{
				Name:        "db_servers",
				Description: "The list of Db servers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DbServers"),
			},
			{
				Name:        "scan_ip_ids",
				Description: "The OCID of the Single Client Access Name (SCAN) IP addresses associated with the cloud VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ScanIpIds"),
			},
			{
				Name:        "vip_ids",
				Description: "The OCID of the virtual IP (VIP) addresses associated with the cloud VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("VipIds"),
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("FreeformTags"),
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DefinedTags"),
			},
			{
				Name:        "data_collection_options",
				Description: "The data collection options of the cloud VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DataCollectionOptions"),
			},
			{
				Name:        "iorm_config_cache",
				Description: "The config cache details of the cloud VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Iormconfigcache"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(cloudVMClusterTags),
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

func listDatabaseCloudVMClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != d.EqualsQualString("compartment_id") {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		logger.Error("oci_database_cloud_vm_cluster.listDatabaseCloudVMClusters", "session_error", err)
		return nil, err
	}

	request := database.ListCloudVmClustersRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["display_name"] != nil {
		dbName := d.EqualsQualString("display_name")
		request.DisplayName = types.String(dbName)
	}

	if equalQuals["cloud_exadata_infrastructure_id"] != nil {
		dbName := d.EqualsQualString("cloud_exadata_infrastructure_id")
		request.CloudExadataInfrastructureId = types.String(dbName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := d.EqualsQualString("lifecycle_state")
		if isValidCloudVMClusterSummaryLifecycleState(lifecycleState) {
			request.LifecycleState = database.CloudVmClusterSummaryLifecycleStateEnum(lifecycleState)
		} else {
			return nil, nil
		}
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DatabaseClient.ListCloudVmClusters(ctx, request)
		if err != nil {
			logger.Error("oci_database_cloud_vm_cluster.listDatabaseCloudVMClusters", "api_error", err)
			return nil, err
		}

		for _, vm := range response.Items {
			d.StreamListItem(ctx, vm)

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

func getDatabaseCloudVmCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d. EqualsQualString("id")

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		logger.Error("oci_database_cloud_vm_cluster.getDatabaseCloudVmCluster", "session_error", err)
		return nil, err
	}

	request := database.GetCloudVmClusterRequest{
		CloudVmClusterId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetCloudVmCluster(ctx, request)
	if err != nil {
		logger.Error("oci_database_cloud_vm_cluster.getDatabaseCloudVmCluster", "api_error", err)
		return nil, err
	}
	return response.CloudVmCluster, nil
}

//// TRANSFORM FUNCTION

func cloudVMClusterTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.CloudVmClusterSummary:
		vmClusterSummary := d.HydrateItem.(database.CloudVmClusterSummary)
		freeformTags = vmClusterSummary.FreeformTags
		definedTags = vmClusterSummary.DefinedTags
	case database.CloudVmCluster:
		pDatabase := d.HydrateItem.(database.CloudVmCluster)
		freeformTags = pDatabase.FreeformTags
		definedTags = pDatabase.DefinedTags
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

func isValidCloudVMClusterSummaryLifecycleState(state string) bool {
	stateType := database.CloudVmClusterLifecycleStateEnum(state)
	switch stateType {
	case database.CloudVmClusterLifecycleStateProvisioning, database.CloudVmClusterLifecycleStateAvailable, database.CloudVmClusterLifecycleStateUpdating, database.CloudVmClusterLifecycleStateTerminating, database.CloudVmClusterLifecycleStateTerminated, database.CloudVmClusterLifecycleStateFailed, database.CloudVmClusterLifecycleStateMaintenanceInProgress:
		return true
	}
	return false
}
