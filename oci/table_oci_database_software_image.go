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

func tableOciDatabaseSoftwareImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_database_software_image",
		Description: "OCI Database Software Image",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSoftwareImage,
		},
		List: &plugin.ListConfig{
			Hydrate: listSoftwareImages,
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
					Name:    "image_type",
					Require: plugin.Optional,
				},
				{
					Name:    "image_shape_family",
					Require: plugin.Optional,
				},
				{
					Name:    "is_upgrade_supported",
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
				Description: "The user-friendly name for the database software image. The name does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the database software image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "image_type",
				Description: "The type of software image. It can be grid or database.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the database software image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the database software image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "database_version",
				Description: "The database version with which the database software image is to be built.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_shape_family",
				Description: "The shape that the image is meant for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "included_patches_summary",
				Description: "The patches included in the image and the version of the image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Detailed message for the lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ls_inventory",
				Description: "The output from lsinventory which will get passed as a string.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "patch_set",
				Description: "The PSU or PBP or release updates.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_upgrade_supported",
				Description: "True if this database software image is supported for upgrade.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "database_software_image_included_patches",
				Description: "List of one-off patches for database homes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "database_software_image_one_off_patches",
				Description: "List of one-off patches for database homes.",
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
				Transform:   transform.From(softwareImageTags),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listSoftwareImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listSoftwareImages", "Compartment", compartment, "OCI_REGION", region)

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

	// Build request parameters
	request := buildDatabaseSoftwareImageFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
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
		response, err := session.DatabaseClient.ListDatabaseSoftwareImages(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, image := range response.Items {
			d.StreamListItem(ctx, image)

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

func getSoftwareImage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getSoftwareImage", "Compartment", compartment, "OCI_REGION", region)

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

	request := database.GetDatabaseSoftwareImageRequest{
		DatabaseSoftwareImageId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.DatabaseClient.GetDatabaseSoftwareImage(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DatabaseSoftwareImage, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func softwareImageTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case database.DatabaseSoftwareImage:
		image := d.HydrateItem.(database.DatabaseSoftwareImage)
		freeformTags = image.FreeformTags
		definedTags = image.DefinedTags
	case database.DatabaseSoftwareImageSummary:
		image := d.HydrateItem.(database.DatabaseSoftwareImageSummary)
		freeformTags = image.FreeformTags
		definedTags = image.DefinedTags
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
func buildDatabaseSoftwareImageFilters(equalQuals plugin.KeyColumnEqualsQualMap) database.ListDatabaseSoftwareImagesRequest {
	request := database.ListDatabaseSoftwareImagesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["image_type"] != nil {
		request.ImageType = database.DatabaseSoftwareImageSummaryImageTypeEnum(equalQuals["image_type"].GetStringValue())
	}
	if equalQuals["image_shape_family"] != nil {
		request.ImageShapeFamily = database.DatabaseSoftwareImageSummaryImageShapeFamilyEnum(equalQuals["image_shape_family"].GetStringValue())
	}
	if equalQuals["is_upgrade_supported"] != nil {
		request.IsUpgradeSupported = types.Bool(equalQuals["is_upgrade_supported"].GetBoolValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = database.DatabaseSoftwareImageSummaryLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
