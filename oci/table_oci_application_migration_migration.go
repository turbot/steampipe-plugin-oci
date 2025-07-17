package oci

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// // TABLE DEFINITION
func tableApplicationMigrationMigration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_application_migration_migration",
		Description: "[DEPRECATED] OCI Application Migration Migration",
		List: &plugin.ListConfig{
			Hydrate: listApplicationMigrationMigrations,
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
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "User-friendly name of the migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description of the migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source_id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the source with which this migration is associated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_name",
				Description: "Name of the application which is being migrated. This is the name of the application in the source environment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "application_type",
				Description: "The type of application being migrated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "pre_created_target_database_type",
				Description: "The pre-existing database type to be used in this migration. Currently, Application migration only supports Oracle Cloud.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_selective_migration",
				Description: "If set to `true`, Application Migration migrates only the application resources that you specify. If set to `false`, Application Migration migrates the entire application. When you migrate the entire application, all the application resources are migrated to the target environment. You can selectively migrate resources only for the Oracle Integration Cloud and Oracle Integration Cloud Service applications.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "service_config",
				Description: "Configuration required to migrate the application. In addition to the key and value, additional fields are provided.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "application_config",
				Description: "Configuration required to migrate the application. In addition to the key and value, additional fields are provided.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Details about the current lifecycle state of the migration.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "migration_state",
				Description: "The current state of the overall migration process.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a namespace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_created",
				Description: "Time that Migration was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
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
		},
	}
}

//// LIST FUNCTION

func listApplicationMigrationMigrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	err := errors.New("The oci_application_migration_migration table has been deprecated and removed, please use oci_cloud_migrations_migration table instead.")
	return nil, err
}
