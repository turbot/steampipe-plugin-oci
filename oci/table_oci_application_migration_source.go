package oci

import (
		"context"
		"github.com/oracle/oci-go-sdk/v65/common"
		"github.com/oracle/oci-go-sdk/v65/applicationmigration"
		"github.com/turbot/go-kit/types"
		"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
		"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
		"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
		"strings"
)

//// TABLE DEFINITION
func tableApplicationMigrationSource(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_application_migration_source",
		Description:      "OCI Source",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getApplicationMigrationSource,
		},
		List: &plugin.ListConfig{
			Hydrate: listApplicationMigrationSources,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:		"compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:		"id",
					Require: plugin.Optional,
				},
				{
					Name:		"display_name",
					Require: plugin.Optional,
				},
				{
					Name: 	"lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "Name of the source. This helps you to identify the appropriate source environment when you have multiple sources defined.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "Description of the source. This helps you to identify the appropriate source environment when you have multiple sources defined.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Details about the current lifecycle state of the source.",
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
				Description: "Time that Source was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

  		// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(applicationMigrationSourceTags),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION
func listApplicationMigrationSources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listApplicationMigrationSources", "Compartment", compartment, "OCI_REGION", region)


	equalQuals := d.KeyColumnQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := applicationMigrationService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildApplicationMigrationSourceFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(100)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ApplicationMigrationClient.ListSources(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
func getApplicationMigrationSource(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getApplicationMigrationSource", "Compartment", compartment, "OCI_REGION", region)



	var id string
	if h.Item != nil {
		id = *h.Item.(applicationmigration.SourceSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := applicationMigrationService(ctx, d, region)
	if err != nil {
		logger.Error("getApplicationMigrationSource", "error_ApplicationMigrationService", err)
		return nil, err
	}

	request := applicationmigration.GetSourceRequest{
		SourceId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ApplicationMigrationClient.GetSource(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Source, nil
}

//// TRANSFORM FUNCTION
func applicationMigrationSourceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case applicationmigration.Source:
		obj := d.HydrateItem.(applicationmigration.Source)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case applicationmigration.SourceSummary:
		obj := d.HydrateItem.(applicationmigration.SourceSummary)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
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

// Build additional filters
func buildApplicationMigrationSourceFilters(equalQuals plugin.KeyColumnEqualsQualMap) applicationmigration.ListSourcesRequest {
	request := applicationmigration.ListSourcesRequest{}

	if equalQuals["compartment_id"] != nil {
			request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())

	}
	if equalQuals["id"] != nil {
			request.Id = types.String(equalQuals["id"].GetStringValue())

	}
	if equalQuals["display_name"] != nil {
			request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())

	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = applicationmigration.ListSourcesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}