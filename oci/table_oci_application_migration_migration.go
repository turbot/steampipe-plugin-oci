package oci

// import (
// 	"context"
// 	"strings"

// 	"github.com/oracle/oci-go-sdk/v65/applicationmigration"
// 	"github.com/oracle/oci-go-sdk/v65/common"
// 	"github.com/turbot/go-kit/types"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
// )

// // // TABLE DEFINITION
// func tableApplicationMigrationMigration(_ context.Context) *plugin.Table {
// 	return &plugin.Table{
// 		Name:             "oci_application_migration_migration",
// 		Description:      "OCI Application Migration Migration",
// 		DefaultTransform: transform.FromCamel(),
// 		Get: &plugin.GetConfig{
// 			KeyColumns: plugin.SingleColumn("id"),
// 			Hydrate:    getApplicationMigrationMigration,
// 		},
// 		List: &plugin.ListConfig{
// 			Hydrate: listApplicationMigrationMigrations,
// 			KeyColumns: []*plugin.KeyColumn{
// 				{
// 					Name:    "compartment_id",
// 					Require: plugin.Optional,
// 				},
// 				{
// 					Name:    "display_name",
// 					Require: plugin.Optional,
// 				},
// 				{
// 					Name:    "lifecycle_state",
// 					Require: plugin.Optional,
// 				},
// 			},
// 		},
// 		GetMatrixItemFunc: BuildCompartementRegionList,
// 		Columns: []*plugin.Column{
// 			{
// 				Name:        "id",
// 				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the migration.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "display_name",
// 				Description: "User-friendly name of the migration.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "description",
// 				Description: "Description of the migration.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "source_id",
// 				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the source with which this migration is associated.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "application_name",
// 				Description: "Name of the application which is being migrated. This is the name of the application in the source environment.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "application_type",
// 				Description: "The type of application being migrated.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "pre_created_target_database_type",
// 				Description: "The pre-existing database type to be used in this migration. Currently, Application migration only supports Oracle Cloud.",
// 				Type:        proto.ColumnType_STRING,
// 				Hydrate:     getApplicationMigrationMigration,
// 			},
// 			{
// 				Name:        "is_selective_migration",
// 				Description: "If set to `true`, Application Migration migrates only the application resources that you specify. If set to `false`, Application Migration migrates the entire application. When you migrate the entire application, all the application resources are migrated to the target environment. You can selectively migrate resources only for the Oracle Integration Cloud and Oracle Integration Cloud Service applications.",
// 				Type:        proto.ColumnType_BOOL,
// 				Hydrate:     getApplicationMigrationMigration,
// 			},
// 			{
// 				Name:        "service_config",
// 				Description: "Configuration required to migrate the application. In addition to the key and value, additional fields are provided.",
// 				Type:        proto.ColumnType_JSON,
// 				Hydrate:     getApplicationMigrationMigration,
// 			},
// 			{
// 				Name:        "application_config",
// 				Description: "Configuration required to migrate the application. In addition to the key and value, additional fields are provided.",
// 				Type:        proto.ColumnType_JSON,
// 				Hydrate:     getApplicationMigrationMigration,
// 			},
// 			{
// 				Name:        "lifecycle_state",
// 				Description: "The current state of the migration.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "lifecycle_details",
// 				Description: "Details about the current lifecycle state of the migration.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "migration_state",
// 				Description: "The current state of the overall migration process.",
// 				Type:        proto.ColumnType_STRING,
// 			},
// 			{
// 				Name:        "freeform_tags",
// 				Description: "Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.",
// 				Type:        proto.ColumnType_JSON,
// 			},
// 			{
// 				Name:        "defined_tags",
// 				Description: "Defined tags for this resource. Each key is predefined and scoped to a namespace.",
// 				Type:        proto.ColumnType_JSON,
// 			},
// 			{
// 				Name:        "time_created",
// 				Description: "Time that Migration was created.",
// 				Type:        proto.ColumnType_TIMESTAMP,
// 				Transform:   transform.FromField("TimeCreated.Time"),
// 			},

// 			// Standard Steampipe columns
// 			{
// 				Name:        "tags",
// 				Description: ColumnDescriptionTags,
// 				Type:        proto.ColumnType_JSON,
// 				Transform:   transform.From(applicationMigrationMigrationTags),
// 			},
// 			{
// 				Name:        "title",
// 				Description: ColumnDescriptionTitle,
// 				Type:        proto.ColumnType_STRING,
// 				Transform:   transform.FromField("DisplayName"),
// 			},

// 			// Standard OCI columns
// 			{
// 				Name:        "compartment_id",
// 				Description: ColumnDescriptionCompartment,
// 				Type:        proto.ColumnType_STRING,
// 				Transform:   transform.FromField("CompartmentId"),
// 			},
// 			{
// 				Name:        "tenant_id",
// 				Description: ColumnDescriptionTenantId,
// 				Type:        proto.ColumnType_STRING,
// 				Hydrate:     getTenantId,
// 				Transform:   transform.FromValue(),
// 			},
// 		},
// 	}
// }

// //// LIST FUNCTION

// func listApplicationMigrationMigrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	region := d.EqualsQualString(matrixKeyRegion)
// 	compartment := d.EqualsQualString(matrixKeyCompartment)
// 	logger.Debug("oci_application_migration_migration.listApplicationMigrationMigrations", "Compartment", compartment, "OCI_REGION", region)

// 	equalQuals := d.EqualsQuals
// 	// Return nil, if given compartment_id doesn't match
// 	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
// 		return nil, nil
// 	}
// 	// Create Session
// 	session, err := applicationMigrationService(ctx, d, region)
// 	if err != nil {
// 		logger.Error("oci_application_migration_migration.listApplicationMigrationMigrations", "connection_error", err)
// 		return nil, err
// 	}

// 	//Build request parameters
// 	request := buildApplicationMigrationMigrationFilters(equalQuals)
// 	request.CompartmentId = types.String(compartment)
// 	request.Limit = types.Int(100)
// 	request.RequestMetadata = common.RequestMetadata{
// 		RetryPolicy: getDefaultRetryPolicy(d.Connection),
// 	}

// 	limit := d.QueryContext.Limit
// 	if d.QueryContext.Limit != nil {
// 		if *limit < int64(*request.Limit) {
// 			request.Limit = types.Int(int(*limit))
// 		}
// 	}

// 	pagesLeft := true
// 	for pagesLeft {
// 		response, err := session.ApplicationMigrationClient.ListMigrations(ctx, request)
// 		if err != nil {
// 			logger.Error("oci_application_migration_migration.listApplicationMigrationMigrations", "api_error", err)
// 			return nil, err
// 		}
// 		for _, respItem := range response.Items {
// 			d.StreamListItem(ctx, respItem)

// 			// Context can be cancelled due to manual cancellation or the limit has been hit
// 			if d.RowsRemaining(ctx) == 0 {
// 				return nil, nil
// 			}
// 		}
// 		if response.OpcNextPage != nil {
// 			request.Page = response.OpcNextPage
// 		} else {
// 			pagesLeft = false
// 		}
// 	}

// 	return nil, err
// }

// //// HYDRATE FUNCTION

// func getApplicationMigrationMigration(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	region := d.EqualsQualString(matrixKeyRegion)
// 	compartment := d.EqualsQualString(matrixKeyCompartment)
// 	logger.Debug("oci_application_migration_migration.getApplicationMigrationMigration", "Compartment", compartment, "OCI_REGION", region)

// 	var id string
// 	if h.Item != nil {
// 		id = *h.Item.(applicationmigration.MigrationSummary).Id
// 	} else {
// 		id = d.EqualsQuals["id"].GetStringValue()
// 		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
// 			return nil, nil
// 		}
// 	}

// 	// handle empty id in get call
// 	if id == "" {
// 		return nil, nil
// 	}

// 	// Create Session

// 	session, err := applicationMigrationService(ctx, d, region)
// 	if err != nil {
// 		logger.Error("oci_application_migration_migration.getApplicationMigrationMigration", "connection_error", err)
// 		return nil, err
// 	}

// 	request := applicationmigration.GetMigrationRequest{
// 		MigrationId: types.String(id),
// 		RequestMetadata: common.RequestMetadata{
// 			RetryPolicy: getDefaultRetryPolicy(d.Connection),
// 		},
// 	}

// 	response, err := session.ApplicationMigrationClient.GetMigration(ctx, request)
// 	if err != nil {
// 		logger.Error("oci_application_migration_migration.getApplicationMigrationMigration", "api_error", err)
// 		return nil, err
// 	}
// 	return response.Migration, nil
// }

// //// TRANSFORM FUNCTION

// func applicationMigrationMigrationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	var freeformTags map[string]string
// 	var definedTags map[string]map[string]interface{}
// 	switch d.HydrateItem.(type) {
// 	case applicationmigration.Migration:
// 		obj := d.HydrateItem.(applicationmigration.Migration)
// 		freeformTags = obj.FreeformTags
// 		definedTags = obj.DefinedTags
// 	case applicationmigration.MigrationSummary:
// 		obj := d.HydrateItem.(applicationmigration.MigrationSummary)
// 		freeformTags = obj.FreeformTags
// 		definedTags = obj.DefinedTags
// 	}

// 	var tags map[string]interface{}
// 	if freeformTags != nil {
// 		tags = map[string]interface{}{}
// 		for k, v := range freeformTags {
// 			tags[k] = v
// 		}
// 	}
// 	if definedTags != nil {
// 		if tags == nil {
// 			tags = map[string]interface{}{}
// 		}
// 		for _, v := range definedTags {
// 			for key, value := range v {
// 				tags[key] = value
// 			}

// 		}
// 	}
// 	return tags, nil
// }

// // Build additional filters
// func buildApplicationMigrationMigrationFilters(equalQuals plugin.KeyColumnEqualsQualMap) applicationmigration.ListMigrationsRequest {
// 	request := applicationmigration.ListMigrationsRequest{}

// 	if equalQuals["display_name"] != nil {
// 		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())

// 	}
// 	if equalQuals["lifecycle_state"] != nil {
// 		request.LifecycleState = applicationmigration.ListMigrationsLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
// 	}

// 	return request
// }
