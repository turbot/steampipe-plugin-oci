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

func tableOciDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_db",
		Description: "OCI Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getDatabase,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			ParentHydrate:     listDatabaseDBHomes,
			Hydrate:           listDatabases,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "db_name",
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
				Name:        "db_name",
				Description: "The database name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "db_unique_name",
				Description: "A system-generated name for the database to ensure uniqueness within an oracle data guard group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the database was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "character_set",
				Description: "The character set for the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_software_image_id",
				Description: "The database software image OCID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "kms_key_version_id",
				Description: "The OCID of the key container version used in TDE operations.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the Oracle Cloud Infrastructure vault.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sid_prefix",
				Description: "Specifies a prefix for the Oracle SID of the database to be created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_cdb",
				Description: "Specifies a prefix for the Oracle SID of the database to be created.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "db_home_id",
				Description: "The OCID of the database home.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "db_system_id",
				Description: "The OCID of the DB system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "db_workload",
				Description: "The database workload type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the key container that is used as the master encryption key in database transparent data encryption (TDE) operations.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "last_backup_timestamp",
				Description: "The date and time when the latest database backup was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LastBackupTimestamp.Time"),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ncharacter_set",
				Description: "The national character set for the database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pdb_name",
				Description: "The name of the pluggable database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_database_point_in_time_recovery_timestamp",
				Description: "Point in time recovery timeStamp of the source database at which cloned database system is cloned from the source database system.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("SourceDatabasePointInTimeRecoveryTimestamp.Time"),
			},
			{
				Name:        "vm_cluster_id",
				Description: "The OCID of the vm cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "connection_strings",
				Description: "The connection strings used to connect to the oracle database.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "db_backup_config",
				Description: "Database backup configuration details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "database_management_config",
				Description: "The configuration of the Database Management service.",
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
				Transform:   transform.From(databaseTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DbName"),
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

func listDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listDatabases", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	homeId := h.Item.(database.DbHomeSummary).Id

	// Create Session
	session, err := databaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := database.ListDatabasesRequest{
		CompartmentId: types.String(compartment),
		DbHomeId:      homeId,
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["db_name"] != nil {
		dbName := equalQuals["db_name"].GetStringValue()
		request.DbName = types.String(dbName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = database.DatabaseSummaryLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.DatabaseClient.ListDatabases(ctx, request)
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

func getDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getDatabase", "Compartment", compartment, "OCI_REGION", region)

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

	request := database.GetDatabaseRequest{
		DatabaseId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetDatabase(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Database, nil
}

//// TRANSFORM FUNCTION

func databaseTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.DatabaseSummary:
		database := d.HydrateItem.(database.DatabaseSummary)
		freeformTags = database.FreeformTags
		definedTags = database.DefinedTags
	case database.Database:
		database := d.HydrateItem.(database.Database)
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
