package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/cloudguard"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudGuardDetectorRecipe(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_detector_recipe",
		Description: "OCI Cloud Guard Detector Recipe",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudGuardDetectorRecipe,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardDetectorRecipes,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartmentList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "name",
				Description: "DisplayName of detector recipe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "Ocid for detector recipe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "source_detector_recipe_id",
				Description: "Recipe Ocid of the Source Recipe to be cloned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the detector recipe.",
				Type:        proto.ColumnType_STRING,
			},

			// other columns
			{
				Name:        "description",
				Description: "Detector recipe description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the detector recipe was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the detector recipe was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "owner",
				Description: "Owner of detector recipe.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "detector",
				Description: "Type of detector.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "detector_rules",
				Description: "List of detector rules for the detector type for recipe.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "effective_detector_rules",
				Description: "List of detector rules for the detector type for recipe.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudGuardDetectorRecipe,
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(cloudGuardDetectorRecipeTags),
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
		}),
	}
}

//// LIST FUNCTION

func listCloudGuardDetectorRecipes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Trace("oci.listCloudGuardDetectorRecipes", "Compartment", compartment)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// fetch reporting region from configuration
	getCloudGuardConfigurationCached := plugin.HydrateFunc(getCloudGuardConfiguration).WithCache()
	configuration, err := getCloudGuardConfigurationCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	reportingRegion := configuration.(cloudguard.Configuration).ReportingRegion

	// Create Session
	session, err := cloudGuardService(ctx, d, *reportingRegion)
	if err != nil {
		return nil, err
	}

	request := cloudguard.ListDetectorRecipesRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["name"] != nil {
		displayName := equalQuals["name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = cloudguard.ListDetectorRecipesLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudGuardClient.ListDetectorRecipes(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, detectorRecipe := range response.Items {
			d.StreamListItem(ctx, detectorRecipe)

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

func getCloudGuardDetectorRecipe(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.getCloudGuardDetectorRecipe", "Compartment", compartment)

	var id string
	if h.Item != nil {
		id = *h.Item.(cloudguard.DetectorRecipeSummary).Id
	} else {
		// Restrict the api call to only root compartment
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// fetch reporting region from configuration
	getCloudGuardConfigurationCached := plugin.HydrateFunc(getCloudGuardConfiguration).WithCache()
	configuration, err := getCloudGuardConfigurationCached(ctx, d, h)
	if err != nil {
		return nil, err
	}

	reportingRegion := configuration.(cloudguard.Configuration).ReportingRegion

	// Create Session
	session, err := cloudGuardService(ctx, d, *reportingRegion)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetDetectorRecipeRequest{
		DetectorRecipeId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.CloudGuardClient.GetDetectorRecipe(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.DetectorRecipe, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func cloudGuardDetectorRecipeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case cloudguard.DetectorRecipeSummary:
		detectorRecipe := d.HydrateItem.(cloudguard.DetectorRecipeSummary)
		freeformTags = detectorRecipe.FreeformTags
		definedTags = detectorRecipe.DefinedTags
		systemTags = detectorRecipe.SystemTags
	case cloudguard.DetectorRecipe:
		detectorRecipe := d.HydrateItem.(cloudguard.DetectorRecipe)
		freeformTags = detectorRecipe.FreeformTags
		definedTags = detectorRecipe.DefinedTags
		systemTags = detectorRecipe.SystemTags
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
