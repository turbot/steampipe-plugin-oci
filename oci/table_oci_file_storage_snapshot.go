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

func tableFileStorageSnapshot(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_file_storage_snapshot",
		Description: "OCI File Storage Snapshot",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"400"}),
			Hydrate:           getFileStorageSnapshot,
		},
		List: &plugin.ListConfig{
			Hydrate:       listFileStorageSnapshots,
			ParentHydrate: listFileStorageFileSystems,
		},
		GetMatrixItem: BuildCompartementZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Name of the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the snapshot.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the snapshot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the snapshot was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "file_system_id",
				Description: "The OCID of the file system from which the snapshot was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("FileSystemId"),
			},
			{
				Name:        "provenance_id",
				Description: "An OCID identifying the parent from which this snapshot was cloned.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ProvenanceId"),
			},
			{
				Name:        "is_clone_source",
				Description: "Specifies whether the snapshot has been cloned.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current 'lifecycleState'.",
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

			//  Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(snapshotTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard OCI columns
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

type snapshotInfo struct {
	filestorage.SnapshotSummary
	CompartmentId string
}

//// LIST FUNCTION

func listFileStorageSnapshots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listFileStorageSnapshots", "Compartment", compartment, "zone", zone)

	fileSystem := h.Item.(filestorage.FileSystemSummary)

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.ListSnapshotsRequest{
		FileSystemId: fileSystem.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.FileStorageClient.ListSnapshots(ctx, request)
		if err != nil {
			plugin.Logger(ctx).Trace("GetError", err)
			return nil, err
		}

		for _, snapshots := range response.Items {
			d.StreamLeafListItem(ctx, snapshotInfo{snapshots, compartment})

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

func getFileStorageSnapshot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFileStorageSnapshot")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFileStorageSnapshot", "Compartment", compartment, "OCI_ZONE", zone)

	// Restrict the api call to only root compartment and one zone/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty snapshot id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.GetSnapshotRequest{
		SnapshotId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.FileStorageClient.GetSnapshot(ctx, request)
	if err != nil {
		return nil, err
	}

	snapshot := filestorage.SnapshotSummary{
		FileSystemId:     response.FileSystemId,
		Id:               response.Id,
		Name:             response.Name,
		TimeCreated:      response.TimeCreated,
		LifecycleState:   filestorage.SnapshotSummaryLifecycleStateEnum(response.LifecycleState),
		ProvenanceId:     response.ProvenanceId,
		IsCloneSource:    response.IsCloneSource,
		LifecycleDetails: response.LifecycleDetails,
		FreeformTags:     response.FreeformTags,
		DefinedTags:      response.DefinedTags,
	}
	rowData := snapshotInfo{snapshot, compartment}

	return rowData, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func snapshotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var tags map[string]interface{}

	freeformTags := d.HydrateItem.(snapshotInfo).FreeformTags

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := d.HydrateItem.(snapshotInfo).DefinedTags

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
