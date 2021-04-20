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

func tableCloudGuardManagedList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_cloud_guard_managed_list",
		Description: "OCI Cloud Guard Managed List",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCloudGuardManagedList,
		},
		List: &plugin.ListConfig{
			Hydrate: listCloudGuardManagedLists,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "ManagedList display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "OCID for managed list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "source_managed_list_id",
				Description: "OCID of the source managed list.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the managed list.",
				Type:        proto.ColumnType_STRING,
			},

			// other columns
			{
				Name:        "description",
				Description: "Managed list description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the managed list was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the managed list was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "list_type",
				Description: "Type of the list.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "feed_provider",
				Description: "Provider of the feed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_editable",
				Description: "If this list is editable or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecyle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},

			// json fields
			{
				Name:        "list_items",
				Description: "List of managed list item.",
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
				Description: "Tags added to instances by the service.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(cloudGuardManagedListTags),
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

func listCloudGuardManagedLists(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.listCloudGuardManagedLists", "Compartment", compartment)

	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.ListManagedListsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CloudGuardClient.ListManagedLists(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, managedList := range response.Items {
			d.StreamListItem(ctx, managedList)
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

func getCloudGuardManagedList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getCloudGuardManagedList", "Compartment", compartment)

	// Rstrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := cloudGuardService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetManagedListRequest{
		ManagedListId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.CloudGuardClient.GetManagedList(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.ManagedList, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func cloudGuardManagedListTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case cloudguard.ManagedListSummary:
		managedList := d.HydrateItem.(cloudguard.ManagedListSummary)
		freeformTags = managedList.FreeformTags
		definedTags = managedList.DefinedTags
		systemTags = managedList.SystemTags
	case cloudguard.ManagedList:
		managedList := d.HydrateItem.(cloudguard.ManagedList)
		freeformTags = managedList.FreeformTags
		definedTags = managedList.DefinedTags
		systemTags = managedList.SystemTags
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
