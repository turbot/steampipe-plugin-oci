package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudGuardResponderRecipe(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_responder_recipe",
		Description: "OCI Cloud Guard Responder Recipe",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudGuardResponderRecipe,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardResponderRecipes,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Display name of responder recipe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "OCID for responder recipe.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "source_responder_recipe_id",
				Description: "Recipe OCID of the source recipe to be cloned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the responder recipe.",
				Type:        proto.ColumnType_STRING,
			},

			// other columns
			{
				Name:        "description",
				Description: "Responder recipe description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the responder recipe was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the responder recipe was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "owner",
				Description: "Owner of responder recipe.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "responder_rules",
				Description: "List of responder rules for the responder type for recipe.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "effective_responder_rules",
				Description: "List of responder rules for the responder type for recipe.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudGuardResponderRecipe,
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
				Description: "Tags added to instances by the service.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(cloudGuardResponderRecipeTags),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCloudGuardResponderRecipes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listCloudGuardResponderRecipes", "Compartment", compartment)

	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.ListResponderRecipesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudGuardClient.ListResponderRecipes(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, responderRecipe := range response.Items {
			d.StreamListItem(ctx, responderRecipe)
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

func getCloudGuardResponderRecipe(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCloudGuardResponderRecipe", "Compartment", compartment)

	var id string
	if h.Item != nil {
		id = *h.Item.(cloudguard.ResponderRecipeSummary).Id
	} else {
		// Rstrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty recipe id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetResponderRecipeRequest{
		ResponderRecipeId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.CloudGuardClient.GetResponderRecipe(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.ResponderRecipe, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func cloudGuardResponderRecipeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case cloudguard.ResponderRecipeSummary:
		ResponderRecipe := d.HydrateItem.(cloudguard.ResponderRecipeSummary)
		freeformTags = ResponderRecipe.FreeformTags
		definedTags = ResponderRecipe.DefinedTags
		systemTags = ResponderRecipe.SystemTags
	case cloudguard.ResponderRecipe:
		ResponderRecipe := d.HydrateItem.(cloudguard.ResponderRecipe)
		freeformTags = ResponderRecipe.FreeformTags
		definedTags = ResponderRecipe.DefinedTags
		systemTags = ResponderRecipe.SystemTags
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
