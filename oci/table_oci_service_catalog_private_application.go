package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/servicecatalog"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOciServiceCatalogPrivateApplication(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_service_catalog_private_application",
		Description: "OCI Service Catalog Private Application",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getServiceCatalogPrivateApplication,
		},
		List: &plugin.ListConfig{
			Hydrate: listServiceCatalogPrivateApplications,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "display_name",
				Description: "The name of the private application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique identifier for the private application in Marketplace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The lifecycle state of the private application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compartment_id",
				Description: "The OCID of the compartment where the private application resides.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "package_type",
				Description: "Type of packages within this private application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Date and time the private application was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "Date and time the private application was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "short_description",
				Description: "A short description of the private application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "long_description",
				Description: "A long description of the private application.",
				Type:        proto.ColumnType_STRING,
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(privateApplicationTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(ociRegionName),
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

func listServiceCatalogPrivateApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	// Return nil, if given compartment_id doesn't match
	if d.EqualsQuals["compartment_id"] != nil && compartment != d.EqualsQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := serviceCatalogService(ctx, d, region)
	if err != nil {
		logger.Error("oci_service_catalog_private_application.listServiceCatalogPrivateApplications", "connection_error", err)
		return nil, err
	}

	// Build request parameters
	request := servicecatalog.ListPrivateApplicationsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if d.EqualsQuals["display_name"] != nil {
		displayName := d.EqualsQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	// Pagination handling
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ServiceCatalogClient.ListPrivateApplications(ctx, request)
		if err != nil {
			logger.Error("oci_service_catalog_private_application.listServiceCatalogPrivateApplications", "api_error", err)
			return nil, err
		}

		for _, privateApp := range response.Items {
			d.StreamListItem(ctx, privateApp)

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

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getServiceCatalogPrivateApplication(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
	session, err := serviceCatalogService(ctx, d, region)
	if err != nil {
		logger.Error("oci_service_catalog_private_application.getServiceCatalogPrivateApplication", "connection_error", err)
		return nil, err
	}

	request := servicecatalog.GetPrivateApplicationRequest{
		PrivateApplicationId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ServiceCatalogClient.GetPrivateApplication(ctx, request)
	if err != nil {
		logger.Error("oci_service_catalog_private_application.getServiceCatalogPrivateApplication", "api_error", err)
		return nil, err
	}

	return response.PrivateApplication, nil
}

//// TRANSFORM FUNCTIONS

func privateApplicationTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch item := d.HydrateItem.(type) {
	case servicecatalog.PrivateApplicationSummary:
		// PrivateApplicationSummary doesn't contain tags
		return nil, nil
	case servicecatalog.PrivateApplication:
		freeformTags = item.FreeformTags
		definedTags = item.DefinedTags
	}

	return extractTags(freeformTags, definedTags), nil
}
