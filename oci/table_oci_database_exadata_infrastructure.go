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

func tableOciDatabaseExadataInfrastructure(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_exadata_infrastructure",
		Description: "OCI Database Exadata Infrastructure",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabaseExadataInfrastructure,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listDatabaseExadataInfrastructures,
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
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "The user-friendly name for the Exadata Cloud@Customer infrastructure. The name does not need to be unique.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the Exadata infrastructure.",
			},
			{
				Name:        "lifecycle_state",
				Type:        proto.ColumnType_STRING,
				Description: "The current lifecycle state of the Exadata infrastructure.",
			},
			{
				Name:        "shape",
				Type:        proto.ColumnType_STRING,
				Description: "The shape of the Exadata infrastructure. The shape determines the amount of CPU, storage, and memory resources allocated to the instance.",
			},
			{
				Name:        "time_zone",
				Type:        proto.ColumnType_STRING,
				Description: "The time zone of the Exadata infrastructure.",
			},
			{
				Name:        "cpus_enabled",
				Type:        proto.ColumnType_INT,
				Description: "The number of enabled CPU cores.",
			},
			{
				Name:        "max_cpu_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of CPU cores available.",
			},
			{
				Name:        "memory_size_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The memory allocated in GBs.",
			},
			{
				Name:        "max_memory_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The total memory available in GBs.",
			},
			{
				Name:        "db_node_storage_size_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The local node storage allocated in GBs.",
			},
			{
				Name:        "max_db_node_storage_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The total local node storage available in GBs.",
			},
			{
				Name:        "data_storage_size_in_tbs",
				Type:        proto.ColumnType_DOUBLE,
				Description: "Size, in terabytes, of the DATA disk group.",
			},
			{
				Name:        "max_data_storage_in_tbs",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The total available DATA disk group size.",
			},
			{
				Name:        "rack_serial_number",
				Type:        proto.ColumnType_STRING,
				Description: "The serial number for the Exadata infrastructure.",
			},
			{
				Name:        "storage_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of Exadata storage servers for the Exadata infrastructure.",
			},
			{
				Name:        "additional_storage_count",
				Type:        proto.ColumnType_INT,
				Description: "The requested number of additional storage servers for the Exadata infrastructure.",
			},
			{
				Name:        "activated_storage_count",
				Type:        proto.ColumnType_INT,
				Description: "The requested number of additional storage servers activated for the Exadata infrastructure.",
			},
			{
				Name:        "compute_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of compute servers for the Exadata infrastructure.",
			},
			{
				Name:        "additional_compute_count",
				Type:        proto.ColumnType_INT,
				Description: "The requested number of additional compute servers for the Exadata infrastructure.",
			},
			{
				Name:        "additional_compute_system_model",
				Type:        proto.ColumnType_STRING,
				Description: "Oracle Exadata System Model specification. The system model determines the amount of compute or storage server resources available for use.",
			},
			{
				Name:        "cloud_control_plane_server1",
				Type:        proto.ColumnType_STRING,
				Description: "The IP address for the first control plane server.",
				Transform:   transform.FromField("CloudControlPlaneServer1"),
			},
			{
				Name:        "cloud_control_plane_server2",
				Type:        proto.ColumnType_STRING,
				Description: "The IP address for the second control plane server.",
				Transform:   transform.FromField("CloudControlPlaneServer2"),
			},
			{
				Name:        "netmask",
				Type:        proto.ColumnType_STRING,
				Description: "The netmask for the control plane network.",
			},
			{
				Name:        "gateway",
				Type:        proto.ColumnType_STRING,
				Description: "The gateway for the control plane network.",
			},
			{
				Name:        "admin_network_cidr",
				Type:        proto.ColumnType_STRING,
				Description: "The CIDR block for the Exadata administration network.",
				Transform:   transform.FromField("AdminNetworkCIDR"),
			},
			{
				Name:        "infini_band_network_cidr",
				Type:        proto.ColumnType_STRING,
				Description: "The CIDR block for the Exadata InfiniBand interconnect.",
				Transform:   transform.FromField("InfiniBandNetworkCIDR"),
			},
			{
				Name:        "corporate_proxy",
				Type:        proto.ColumnType_STRING,
				Description: "The corporate network proxy for access to the control plane network.",
			},
			{
				Name:        "time_created",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time the Exadata infrastructure was created.",
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_details",
				Type:        proto.ColumnType_STRING,
				Description: "Additional information about the current lifecycle state.",
			},
			{
				Name:        "csi_number",
				Type:        proto.ColumnType_STRING,
				Description: "The CSI Number of the Exadata infrastructure.",
			},
			{
				Name:        "maintenance_slo_status",
				Type:        proto.ColumnType_STRING,
				Description: "A field to capture ‘Maintenance SLO Status’ for the Exadata infrastructure with values ‘OK’, ‘DEGRADED’. Default is ‘OK’ when the infrastructure is provisioned.",
			},
			{
				Name:        "storage_server_version",
				Type:        proto.ColumnType_STRING,
				Description: "The software version of the storage servers (cells) in the Exadata infrastructure.",
			},
			{
				Name:        "db_server_version",
				Type:        proto.ColumnType_STRING,
				Description: "The software version of the database servers (dom0) in the Exadata infrastructure.",
			},
			{
				Name:        "monthly_db_server_version",
				Type:        proto.ColumnType_STRING,
				Description: "The monthly software version of the database servers (dom0) in the Exadata infrastructure.",
			},
			{
				Name:        "last_maintenance_run_id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the last maintenance run.",
			},
			{
				Name:        "next_maintenance_run_id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the next maintenance run.",
			},
			{
				Name:        "is_cps_offline_report_enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether CPS offline diagnostic report is enabled for this Exadata infrastructure. This will allow a customer to quickly check status themselves and fix problems on their end, saving time and frustration for both Oracle and the customer when they find the CPS in a disconnected state. You can enable offline diagnostic report during Exadata infrastructure provisioning. You can also disable or enable it at any time using the UpdateExadatainfrastructure API.",
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
				Name:        "maintenance_window",
				Type:        proto.ColumnType_JSON,
				Description: "Maintenance Window information.",
			},
			{
				Name:        "contacts",
				Type:        proto.ColumnType_JSON,
				Description: "The list of contacts for the Exadata infrastructure.",
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(exadataInfrastructureDatabaseTags),
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

func listDatabaseExadataInfrastructures(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		logger.Error("oci_database_exadata_infrastructure.listDatabaseExadataInfrastructures", "session_error", err)
		return nil, err
	}

	request := database.ListExadataInfrastructuresRequest{
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

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidExadataInfrastructureSummaryLifecycleState(lifecycleState) {
			request.LifecycleState = database.ExadataInfrastructureSummaryLifecycleStateEnum(lifecycleState)
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
		response, err := session.DatabaseClient.ListExadataInfrastructures(ctx, request)
		if err != nil {
			logger.Error("oci_database_exadata_infrastructure.listDatabaseExadataInfrastructures", "api_error", err)
			return nil, err
		}

		for _, infra := range response.Items {
			d.StreamListItem(ctx, infra)

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

func getDatabaseExadataInfrastructure(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
		logger.Error("oci_database_exadata_infrastructure.getDatabaseExadataInfrastructure", "session_error", err)
		return nil, err
	}

	request := database.GetExadataInfrastructureRequest{
		ExadataInfrastructureId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetExadataInfrastructure(ctx, request)
	if err != nil {
		logger.Error("oci_database_exadata_infrastructure.getDatabaseExadataInfrastructure", "api_error", err)
		return nil, err
	}
	return response.ExadataInfrastructure, nil
}

//// TRANSFORM FUNCTION

func exadataInfrastructureDatabaseTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.ExadataInfrastructureSummary:
		pDatabaseSummary := d.HydrateItem.(database.ExadataInfrastructureSummary)
		freeformTags = pDatabaseSummary.FreeformTags
		definedTags = pDatabaseSummary.DefinedTags
	case database.ExadataInfrastructure:
		pDatabase := d.HydrateItem.(database.ExadataInfrastructure)
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

func isValidExadataInfrastructureSummaryLifecycleState(state string) bool {
	stateType := database.ExadataInfrastructureSummaryLifecycleStateEnum(state)
	switch stateType {
	case database.ExadataInfrastructureSummaryLifecycleStateCreating, database.ExadataInfrastructureSummaryLifecycleStateRequiresActivation, database.ExadataInfrastructureSummaryLifecycleStateActivating, database.ExadataInfrastructureSummaryLifecycleStateActive, database.ExadataInfrastructureSummaryLifecycleStateActivationFailed, database.ExadataInfrastructureSummaryLifecycleStateFailed, database.ExadataInfrastructureSummaryLifecycleStateUpdating, database.ExadataInfrastructureSummaryLifecycleStateDeleting, database.ExadataInfrastructureSummaryLifecycleStateDeleted, database.ExadataInfrastructureSummaryLifecycleStateDisconnected, database.ExadataInfrastructureSummaryLifecycleStateMaintenanceInProgress:
		return true
	}
	return false
}
