package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudGuardTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_target",
		Description: "OCI Cloud Guard Target",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudGuardTarget,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardTargets,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Target display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "OCID for target.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "target_resource_id",
				Description: "Resource ID which the target uses to monitor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the target was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "description",
				Description: "The target description.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCloudGuardTarget,
			},
			{
				Name:        "lifecyle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "recipe_count",
				Description: "Total number of recipes attached to target.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "target_resource_type",
				Description: "Possible type of targets(compartment/HCMCloud/ERPCloud).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The date and time the target was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			// json fields
			{
				Name:        "inherited_by_compartments",
				Description: "List of inherited compartments.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudGuardTarget,
			},
			{
				Name:        "target_detector_recipes",
				Description: "List of detector recipes associated with target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudGuardTarget,
			},
			{
				Name:        "target_responder_recipes",
				Description: "List of responder recipes associated with target.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCloudGuardTarget,
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
				Transform:   transform.From(cloudGuardTargetTags),
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

func listCloudGuardTargets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listCloudGuardTargets", "Compartment", compartment)

	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.ListTargetsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudGuardClient.ListTargets(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, target := range response.Items {
			d.StreamListItem(ctx, target)
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

func getCloudGuardTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCloudGuardTarget", "Compartment", compartment)

	var id string
	if h.Item != nil {
		id = *h.Item.(cloudguard.TargetSummary).Id
	} else {
		// Restrict the api call to only root compartment/ per region
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
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetTargetRequest{
		TargetId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.CloudGuardClient.GetTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Target, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func cloudGuardTargetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case cloudguard.TargetSummary:
		target := d.HydrateItem.(cloudguard.TargetSummary)
		freeformTags = target.FreeformTags
		definedTags = target.DefinedTags
		systemTags = target.SystemTags
	case cloudguard.Target:
		target := d.HydrateItem.(cloudguard.Target)
		freeformTags = target.FreeformTags
		definedTags = target.DefinedTags
		systemTags = target.SystemTags
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
