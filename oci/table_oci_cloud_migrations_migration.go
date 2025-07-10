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

func tableCloudMigrationsMigration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_cloud_migrations_migration",
		Description:      "OCI Cloud Migrations Migration.",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			Hydrate: listCloudMigrationsMigrations,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "compartment_id", Require: plugin.Optional},
				{Name: "lifecycle_state", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudMigrationsMigration,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Unique identifier that is immutable on creation.",
			},
			{
				Name:        "display_name",
				Type:        proto.ColumnType_STRING,
				Description: "Migration Identifier that can be renamed.",
			},
			{
				Name:        "lifecycle_state",
				Type:        proto.ColumnType_STRING,
				Description: "The current state of migration.",
			},
			{
				Name:        "lifecycle_details",
				Type:        proto.ColumnType_STRING,
				Description: "A message describing the current state in more detail.",
			},
			{
				Name:        "time_created",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the migration project was created.",
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time when the migration project was updated.",
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "replication_schedule_id",
				Type:        proto.ColumnType_STRING,
				Description: "Replication schedule identifier.",
			},
			{
				Name:        "is_completed",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether migration is marked as completed.",
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

func listCloudMigrationsMigrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_cloud_migrations_migration.listCloudMigrationsMigrations", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	session, err := cloudMigrationsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration.listCloudMigrationsMigrations", "connection_error", err)
		return nil, err
	}

	request := cloudmigrations.ListMigrationsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Handle optional quals
	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = cloudmigrations.MigrationLifecycleStateEnum(lifecycleState)
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
		response, err := session.CloudMigrationsClient.ListMigrations(ctx, request)
		if err != nil {
			logger.Error("oci_cloud_migrations_migration.listCloudMigrationsMigrations", "api_error", err)
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

func getCloudMigrationsMigration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_cloud_migrations_migration.getCloudMigrationsMigration", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(cloudmigrations.MigrationSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}
	if id == "" {
		return nil, nil
	}

	session, err := cloudMigrationsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration.getCloudMigrationsMigration", "connection_error", err)
		return nil, err
	}

	request := cloudmigrations.GetMigrationRequest{
		MigrationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.CloudMigrationsClient.GetMigration(ctx, request)
	if err != nil {
		logger.Error("oci_cloud_migrations_migration.getCloudMigrationsMigration", "api_error", err)
		return nil, err
	}
	return response.Migration, nil
}
