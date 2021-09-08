package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/nosql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableNoSQLTable(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_nosql_table",
		Description: "OCI NoSQL Table",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getNoSQLTable,
		},
		List: &plugin.ListConfig{
			Hydrate: listNoSQLTables,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Immutable human-friendly table name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Unique identifier that is immutable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The state of a table.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the the table was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "ddl_statement",
				Description: "A DDL statement representing the schema.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getNoSQLTable,
			},
			{
				Name:        "is_auto_reclaimable",
				Description: "True if this table can be reclaimed after an idle period.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_of_expiration",
				Description: "If lifecycleState is INACTIVE, indicates when this table will be automatically removed.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfExpiration.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The time the the table was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "schema",
				Description: "The schema of the table.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNoSQLTable,
			},
			{
				Name:        "table_limits",
				Description: "Various limit for the table.",
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
			{
				Name:        "system_tags",
				Description: ColumnDescriptionSystemTags,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(noSQLTableTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listNoSQLTables(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listNoSQLTables", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := noSQLDatabaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := nosql.ListTablesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.NoSQLClient.ListTables(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, table := range response.Items {
			d.StreamListItem(ctx, table)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				response.OpcNextPage = nil
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

//// HYDRATE FUNCTIONS

func getNoSQLTable(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getNoSQLTable", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(nosql.TableSummary).Id
	} else {
		// Restrict the api call to only root compartment
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := noSQLDatabaseService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := nosql.GetTableRequest{
		TableNameOrId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.NoSQLClient.GetTable(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Table, nil
}

//// TRANSFORM FUNCTION

func noSQLTableTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case nosql.TableSummary:
		table := d.HydrateItem.(nosql.TableSummary)
		freeformTags = table.FreeformTags
		definedTags = table.DefinedTags
		systemTags = table.SystemTags
	case nosql.Table:
		table := d.HydrateItem.(nosql.Table)
		freeformTags = table.FreeformTags
		definedTags = table.DefinedTags
		systemTags = table.SystemTags
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

	if systemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range systemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
