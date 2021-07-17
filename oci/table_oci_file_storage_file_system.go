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

func tableFileStorageFileSystem(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_file_storage_file_system",
		Description: "OCI File Storage File System",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"400"}),
			Hydrate:           getFileStorageFileSystem,
		},
		List: &plugin.ListConfig{
			Hydrate: listFileStorageFileSystems,
		},
		GetMatrixItem: BuildCompartementZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name. It does not have to be unique, and it is changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the file system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain the file system is in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the file system was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "is_clone_parent",
				Description: "Specifies whether the file system has been cloned.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_hydrated",
				Description: "Specifies whether the data has finished copying from the source to the clone.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the KMS key used to encrypt the encryption keys associated with this file system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current 'lifecycleState'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metered_bytes",
				Description: "The number of bytes consumed by the file system.",
				Type:        proto.ColumnType_INT,
			},

			//json fields
			{
				Name:        "source_details",
				Description: "Source information for the file system.",
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

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(fileSystemTags),
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

func listFileStorageFileSystems(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listFileStorageFileSystems", "Compartment", compartment, "zone", zone)

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.ListFileSystemsRequest{
		CompartmentId:      types.String(compartment),
		AvailabilityDomain: types.String(zone),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.FileStorageClient.ListFileSystems(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, fileSystems := range response.Items {
			d.StreamListItem(ctx, fileSystems)
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

func getFileStorageFileSystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getFileStorageFileSystem")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFunctionsApplication", "Compartment", compartment, "OCI_ZONE", zone)

	var id string
	if h.Item != nil {
		fileSystem := h.Item.(filestorage.FileSystemSummary)
		id = *fileSystem.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty application id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := filestorage.GetFileSystemRequest{
		FileSystemId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.FileStorageClient.GetFileSystem(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.FileSystem, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func fileSystemTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	freeformTags := fileSystemFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := fileSystemDefinedTags(d.HydrateItem)

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

func fileSystemFreeformTags(item interface{}) map[string]string {
	switch item.(type) {
	case filestorage.FileSystem:
		return item.(filestorage.FileSystem).FreeformTags
	case filestorage.FileSystemSummary:
		return item.(filestorage.FileSystemSummary).FreeformTags
	}
	return nil
}

func fileSystemDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item.(type) {
	case filestorage.FileSystem:
		return item.(filestorage.FileSystem).DefinedTags
	case filestorage.FileSystemSummary:
		return item.(filestorage.FileSystemSummary).DefinedTags
	}
	return nil
}
