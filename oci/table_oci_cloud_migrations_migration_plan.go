package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/cloudmigrations"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudMigrationsMigrationPlan(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_cloud_migrations_migration_plan",
		Description:      "OCI Cloud Migrations Migration Plan.",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			Hydrate: listCloudMigrationsMigrationPlans,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "migration_id", Require: plugin.Optional},
				{Name: "compartment_id", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "lifecycle_state", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudMigrationsMigrationPlan,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique Oracle ID (OCID) that is immutable on creation.",
			},
			{
				Name:        "migration_id",
				Type:        proto.ColumnType_STRING,
				Description: "The OCID of the associated migration.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
			},
			{
				Name:        "lifecycle_state",
				Type:        proto.ColumnType_STRING,
				Description: "The current state of the migration plan.",
			},
			{
				Name:        "lifecycle_details",
				Type:        proto.ColumnType_STRING,
				Description: "A message describing the current state in more detail.",
			},
			{
				Name:        "time_created",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the migration plan was created.",
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the migration plan was updated.",
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "strategies",
				Type:        proto.ColumnType_JSON,
				Description: "List of strategies for the resources to be migrated.",
			},
			{
				Name:        "calculated_limits",
				Type:        proto.ColumnType_JSON,
				Description: "Limits of the resources that are needed for migration.",
			},
			{
				Name:        "target_environments",
				Type:        proto.ColumnType_JSON,
				Description: "List of target environments.",
			},
			{
				Name:        "migration_plan_stats",
				Type:        proto.ColumnType_JSON,
				Description: "Migration plan statistics.",
			},
			{
				Name:        "reference_to_rms_stack",
				Type:        proto.ColumnType_STRING,
				Description: "OCID of the referenced ORM job.",
			},
			{
				Name:        "source_migration_plan_id",
				Type:        proto.ColumnType_STRING,
				Description: "Source migration plan ID to be cloned.",
			},
			{
				Name:        "freeform_tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionFreefromTags,
			},
			{
				Name:        "defined_tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionDefinedTags,
			},
			{
				Name:        "system_tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionSystemTags,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
				Description: ColumnDescriptionTitle,
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

func listCloudMigrationsMigrationPlans(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_cloud_migrations_migration_plan.listCloudMigrationsMigrationPlans", "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	session, err := cloudMigrationsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration_plan.listCloudMigrationsMigrationPlans", "connection_error", err)
		return nil, err
	}

	request := cloudmigrations.ListMigrationPlansRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Handle optional quals
	if equalQuals["migration_id"] != nil {
		migrationId := equalQuals["migration_id"].GetStringValue()
		request.MigrationId = types.String(migrationId)
	}

	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = cloudmigrations.MigrationPlanLifecycleStateEnum(lifecycleState)
	}

	// Handle limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudMigrationsClient.ListMigrationPlans(ctx, request)
		if err != nil {
			logger.Error("oci_cloud_migrations_migration_plan.listCloudMigrationsMigrationPlans", "api_error", err)
			return nil, err
		}
		for _, item := range response.Items {
			d.StreamListItem(ctx, item)
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
	return nil, nil
}

func getCloudMigrationsMigrationPlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	logger.Debug("oci_cloud_migrations_migration_plan.getCloudMigrationsMigrationPlan", "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(cloudmigrations.MigrationPlanSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}
	if id == "" {
		return nil, nil
	}

	session, err := cloudMigrationsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration_plan.getCloudMigrationsMigrationPlan", "connection_error", err)
		return nil, err
	}

	request := cloudmigrations.GetMigrationPlanRequest{
		MigrationPlanId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.CloudMigrationsClient.GetMigrationPlan(ctx, request)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration_plan.getCloudMigrationsMigrationPlan", "api_error", err)
		return nil, err
	}
	return response.MigrationPlan, nil
}
