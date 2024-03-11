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
		Name:        "oci_database_cloud_exadata_infrastructure",
		Description: "OCI Database Cloud Exadata Infrastructure",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabaseCloudExadataInfrastructure,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listDatabaseCloudExadataInfrastructures,
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
				Description: "The user-friendly name for the cloud Exadata infrastructure resource. The name does not need to be unique.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the Cloud Exadata infrastructure.",
			},
			{
				Name:        "lifecycle_state",
				Type:        proto.ColumnType_STRING,
				Description: "The current lifecycle state of the Cloud Exadata infrastructure.",
			},
			{
				Name:        "shape",
				Type:        proto.ColumnType_STRING,
				Description: "The shape of the Cloud Exadata infrastructure. The shape determines the amount of CPU, storage, and memory resources allocated to the instance.",
			},
			{
				Name:        "time_created",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time the Cloud Exadata infrastructure was created.",
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "max_cpu_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of CPU cores available.",
			},
			{
				Name:        "cpu_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of CPU cores allocated.",
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
				Name:        "storage_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of storage servers for the cloud Exadata infrastructure.",
			},
			{
				Name:        "compute_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of compute servers for the Exadata infrastructure.",
			},
			{
				Name:        "lifecycle_details",
				Type:        proto.ColumnType_STRING,
				Description: "Additional information about the current lifecycle state.",
			},
			{
				Name:        "next_maintenance_run_id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the next maintenance run.",
			},
			{
				Name:        "last_maintenance_run_id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the last maintenance run.",
			},
			{
				Name:        "availability_domain",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the availability domain that the cloud Exadata infrastructure resource is located in.",
			},
			{
				Name:        "total_storage_size_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The total storage allocated to the cloud Exadata infrastructure resource, in gigabytes (GB).",
				Transform: transform.FromField("TotalStorageSizeInGBs"),
			},
			{
				Name:        "available_storage_size_in_gbs",
				Type:        proto.ColumnType_INT,
				Description: "The available storage can be allocated to the cloud Exadata infrastructure resource, in gigabytes (GB).",
				Transform: transform.FromField("AvailableStorageSizeInGBs"),
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
				Name:        "customer_contacts",
				Type:        proto.ColumnType_JSON,
				Description: " The list of customer email addresses that receive information from Oracle about the specified OCI Database service resource.",
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(exadataCloudInfrastructureDatabaseTags),
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

func listDatabaseCloudExadataInfrastructures(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		logger.Error("oci_database_cloud_exadata_infrastructure.listDatabaseCloudExadataInfrastructures", "session_error", err)
		return nil, err
	}

	request := database.ListCloudExadataInfrastructuresRequest{
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

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := d.EqualsQualString("lifecycle_state")
		if isValidExadataCloudInfrastructureSummaryLifecycleState(lifecycleState) {
			request.LifecycleState = database.CloudExadataInfrastructureSummaryLifecycleStateEnum(lifecycleState)
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
		response, err := session.DatabaseClient.ListCloudExadataInfrastructures(ctx, request)
		if err != nil {
			logger.Error("oci_database_cloud_exadata_infrastructure.listDatabaseCloudExadataInfrastructures", "api_error", err)
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

func getDatabaseCloudExadataInfrastructure(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQualString("id")

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		logger.Error("oci_database_cloud_exadata_infrastructure.getDatabaseCloudExadataInfrastructure", "session_error", err)
		return nil, err
	}

	request := database.GetCloudExadataInfrastructureRequest{
		CloudExadataInfrastructureId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetCloudExadataInfrastructure(ctx, request)
	if err != nil {
		logger.Error("oci_database_cloud_exadata_infrastructure.getDatabaseCloudExadataInfrastructure", "api_error", err)
		return nil, err
	}
	return response.CloudExadataInfrastructure, nil
}

//// TRANSFORM FUNCTION

func exadataCloudInfrastructureDatabaseTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.CloudExadataInfrastructureSummary:
		pDatabaseSummary := d.HydrateItem.(database.ExadataInfrastructureSummary)
		freeformTags = pDatabaseSummary.FreeformTags
		definedTags = pDatabaseSummary.DefinedTags
	case database.CloudExadataInfrastructure:
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

func isValidExadataCloudInfrastructureSummaryLifecycleState(state string) bool {
	stateType := database.CloudExadataInfrastructureSummaryLifecycleStateEnum(state)
	switch stateType {
	case database.CloudExadataInfrastructureSummaryLifecycleStateProvisioning, database.CloudExadataInfrastructureSummaryLifecycleStateAvailable, database.CloudExadataInfrastructureSummaryLifecycleStateUpdating, database.CloudExadataInfrastructureSummaryLifecycleStateTerminated, database.CloudExadataInfrastructureSummaryLifecycleStateTerminating, database.CloudExadataInfrastructureSummaryLifecycleStateFailed, database.CloudExadataInfrastructureSummaryLifecycleStateMaintenanceInProgress:
		return true
	}
	return false
}