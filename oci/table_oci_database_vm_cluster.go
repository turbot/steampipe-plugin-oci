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

func tableOciDatabaseVMCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_vm_cluster",
		Description: "OCI Database VM Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabaseVmCluster,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listDatabaseVMClusters,
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
					Name:    "exadata_infrastructure_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "last_patch_history_entry_id",
				Description: "The OCID of the last patch history.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LastPatchHistoryEntryId"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LifecycleState"),
			},
			{
				Name:        "display_name",
				Description: "The user-friendly name for the cloud VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "time_created",
				Description: "The date and time that the VM cluster was created.",
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
				Description: "The time zone of the Exadata infrastructure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeZone"),
			},
			{
				Name:        "is_local_backup_enabled",
				Description: "Indicates if database backup on local Exadata storage is configured for the VM cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsLocalBackupEnabled"),
			},
			{
				Name:        "exadata_infrastructure_id",
				Description: "The OCID of the Exadata infrastructure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ExadataInfrastructureId"),
			},
			{
				Name:        "is_sparse_diskgroup_enabled",
				Description: "Indicates if a sparse disk group is configured for the VM cluster.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsSparseDiskgroupEnabled"),
			},
			{
				Name:        "vm_cluster_network_id",
				Description: "The OCID of the VM cluster network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VmClusterNetworkId"),
			},
			{
				Name:        "cpus_enabled",
				Description: "The number of enabled CPU cores.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CpusEnabled"),
			},
			{
				Name:        "ocpus_enabled",
				Description: "The number of enabled OCPU cores.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("OcpusEnabled"),
			},
			{
				Name:        "memory_size_in_gbs",
				Description: "The memory allocated in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MemorySizeInGBs"),
			},
			{
				Name:        "db_node_storage_size_in_gbs",
				Description: "The local node storage allocated in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DbNodeStorageSizeInGBs"),
			},
			{
				Name:        "data_storage_size_in_tbs",
				Description: "Size, in terabytes, of the DATA disk group.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DataStorageSizeInTBs"),
			},
			{
				Name:        "data_storage_size_in_gbs",
				Description: "Size of the DATA disk group in GBs.",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("DataStorageSizeInGBs"),
			},
			{
				Name:        "shape",
				Description: "The shape of the Exadata infrastructure.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Shape"),
			},
			{
				Name:        "gi_version",
				Description: "The Oracle Grid Infrastructure software version for the VM cluster.",
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
				Description: "The Oracle license model that applies to the VM cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LicenseModel"),
			},
			{
				Name:        "ssh_public_keys",
				Description: "The public key portion of one or more key pairs used for SSH access to the VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SshPublicKeys"),
			},
			{
				Name:        "db_servers",
				Description: "The list of Db servers.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DbServers"),
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
				Description: "Data collection options for the VM cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("DataCollectionOptions"),
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(exadataVMClusterTags),
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

func listDatabaseVMClusters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		logger.Error("oci_database_vm_cluster.listDatabaseVMClusters", "session_error", err)
		return nil, err
	}

	request := database.ListVmClustersRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["display_name"] != nil {
		dbName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(dbName)
	}

	if equalQuals["exadata_infrastructure_id"] != nil {
		infraId := equalQuals["exadata_infrastructure_id"].GetStringValue()
		request.ExadataInfrastructureId = types.String(infraId)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidVMClusterSummaryLifecycleState(lifecycleState) {
			request.LifecycleState = database.VmClusterSummaryLifecycleStateEnum(lifecycleState)
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
		response, err := session.DatabaseClient.ListVmClusters(ctx, request)
		if err != nil {
			logger.Error("oci_database_vm_cluster.listDatabaseVMClusters", "api_error", err)
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

func getDatabaseVmCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

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
		logger.Error("oci_database_vm_cluster.getDatabaseVmCluster", "session_error", err)
		return nil, err
	}

	request := database.GetVmClusterRequest{
		VmClusterId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetVmCluster(ctx, request)
	if err != nil {
		logger.Error("oci_database_vm_cluster.getDatabaseVmCluster", "api_error", err)
		return nil, err
	}
	return response.VmCluster, nil
}

//// TRANSFORM FUNCTION

func exadataVMClusterTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.VmClusterSummary:
		vmClusterSummary := d.HydrateItem.(database.VmClusterSummary)
		freeformTags = vmClusterSummary.FreeformTags
		definedTags = vmClusterSummary.DefinedTags
	case database.VmCluster:
		cluster := d.HydrateItem.(database.VmCluster)
		freeformTags = cluster.FreeformTags
		definedTags = cluster.DefinedTags
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

func isValidVMClusterSummaryLifecycleState(state string) bool {
	stateType := database.VmClusterLifecycleStateEnum(state)
	switch stateType {
	case database.VmClusterLifecycleStateProvisioning, database.VmClusterLifecycleStateAvailable, database.VmClusterLifecycleStateUpdating, database.VmClusterLifecycleStateTerminating, database.VmClusterLifecycleStateTerminated, database.VmClusterLifecycleStateFailed, database.VmClusterLifecycleStateMaintenanceInProgress:
		return true
	}
	return false
}
