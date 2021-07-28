package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/filestorage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableFileStorageMountTarget(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_file_storage_mount_target",
		Description: "OCI File Storage Mount Target",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"400"}),
			Hydrate:           getFileStorageMountTarget,
		},
		List: &plugin.ListConfig{
			Hydrate: listFileStorageMountTargets,
		},
		GetMatrixItem: BuildCompartementZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name of the Mount Target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the Mount Target.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Mount Target.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain the mount target is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the Mount Target was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "export_set_id",
				Description: "The OCID of the associated export set.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current 'lifecycleState'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getFileStorageMountTarget,
			},
			{
				Name:        "subnet_id",
				Description: "The OCIDs of the subnet the mount target is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "nsg_ids",
				Description: "A list of Network Security Group OCIDs associated with this mount target.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "private_ip_ids",
				Description: "The OCIDs of the private IP addresses associated with this mount target.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel(),
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

			//  Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(mountTargetTags),
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
				Transform:   transform.FromField("SubnetId").Transform(ociRegionName),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listFileStorageMountTargets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listFileStorageMountTargets", "Compartment", compartment, "zone", zone)

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.ListMountTargetsRequest{
		CompartmentId:      types.String(compartment),
		AvailabilityDomain: types.String(zone),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.FileStorageClient.ListMountTargets(ctx, request)
		if err != nil {
			plugin.Logger(ctx).Trace("GetError", err)
			return nil, err
		}

		for _, mountTarget := range response.Items {
			d.StreamLeafListItem(ctx, mountTarget)
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

func getFileStorageMountTarget(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFileStorageMountTarget")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFileStorageMountTarget", "Compartment", compartment, "OCI_ZONE", zone)

	id := d.KeyColumnQuals["id"].GetStringValue()
	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	// handle empty snapshot id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.GetMountTargetRequest{
		MountTargetId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.FileStorageClient.GetMountTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func mountTargetTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case filestorage.MountTargetSummary:
		mountTarget := d.HydrateItem.(filestorage.MountTargetSummary)
		freeformTags = mountTarget.FreeformTags
		definedTags = mountTarget.DefinedTags
	case filestorage.MountTarget:
		mountTarget := d.HydrateItem.(filestorage.MountTarget)
		freeformTags = mountTarget.FreeformTags
		definedTags = mountTarget.DefinedTags
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
