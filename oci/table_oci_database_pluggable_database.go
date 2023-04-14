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

func tableOciPluggableDatabase(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_pluggable_database",
		Description: "OCI Pluggable Database",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getPluggableDatabase,
		},
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           listDatabasePluggableDatabases,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "pdb_name",
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
				Name:        "pdb_name",
				Description: "The name for the pluggable database. The name is unique in the context of a Database. The name must begin with an alphabetic character and can contain a maximum of thirty alphanumeric characters. Special characters are not permitted. The pluggable database name should not be same as the container database name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the pluggable database.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the pluggable database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the pluggable database was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "container_database_id",
				Description: "The OCID of the CDB.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "is_restricted",
				Description: "The restricted mode of pluggableDatabase. If a pluggableDatabase is opened in restricted mode, the user needs both Create a session and restricted session privileges to connect to it.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_details",
				Description: "Detailed message for the lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "open_mode",
				Description: "The mode that pluggableDatabase is in. Open mode can only be changed to READ_ONLY or MIGRATE directly from the backend.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection_strings",
				Description: "The connection strings used to connect to the oracle pluggable database.",
				Type:        proto.ColumnType_JSON,
			},

			// Tags
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
				Transform:   transform.From(pluggableDatabaseTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("PdbName"),
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

func listDatabasePluggableDatabases(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listPluggableDatabases", "Compartment", compartment, "OCI_REGION", region)

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

	request := database.ListPluggableDatabasesRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["pdb_name"] != nil {
		dbName := equalQuals["pdb_name"].GetStringValue()
		request.PdbName = types.String(dbName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidPluggableDatabaseSummaryLifecycleState(lifecycleState) {
			request.LifecycleState = database.PluggableDatabaseSummaryLifecycleStateEnum(lifecycleState)
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
		response, err := session.DatabaseClient.ListPluggableDatabases(ctx, request)
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

func getPluggableDatabase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getPluggableDatabase", "Compartment", compartment, "OCI_REGION", region)

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

	request := database.GetPluggableDatabaseRequest{
		PluggableDatabaseId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetPluggableDatabase(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.PluggableDatabase, nil
}

//// TRANSFORM FUNCTION

func pluggableDatabaseTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.PluggableDatabaseSummary:
		pDatabaseSummary := d.HydrateItem.(database.PluggableDatabaseSummary)
		freeformTags = pDatabaseSummary.FreeformTags
		definedTags = pDatabaseSummary.DefinedTags
	case database.PluggableDatabase:
		pDatabase := d.HydrateItem.(database.PluggableDatabase)
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

func isValidPluggableDatabaseSummaryLifecycleState(state string) bool {
	stateType := database.PluggableDatabaseSummaryLifecycleStateEnum(state)
	switch stateType {
	case database.PluggableDatabaseSummaryLifecycleStateProvisioning, database.PluggableDatabaseSummaryLifecycleStateAvailable, database.PluggableDatabaseSummaryLifecycleStateFailed, database.PluggableDatabaseSummaryLifecycleStateTerminated, database.PluggableDatabaseSummaryLifecycleStateTerminating, database.PluggableDatabaseSummaryLifecycleStateUpdating:
		return true
	}
	return false
}
