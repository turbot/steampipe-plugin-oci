package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/filestorage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartmentZonalList,
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
				Name:        "exports",
				Description: "A list of export resources by file system.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getFileStorageFileSystemExports,
				Transform:   transform.FromValue(),
			},
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
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

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Return nil, if given availability_domain doesn't match
	if equalQuals["availability_domain"] != nil && zone != equalQuals["availability_domain"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := fileStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildFileStorageFileSystemFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.AvailabilityDomain = types.String(zone)
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
		response, err := session.FileStorageClient.ListFileSystems(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, fileSystems := range response.Items {
			d.StreamListItem(ctx, fileSystems)

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

func getFileStorageFileSystem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getFileStorageFileSystem", "Compartment", compartment, "OCI_ZONE", zone)

	var id string
	if h.Item != nil {
		fileSystem := h.Item.(filestorage.FileSystemSummary)
		id = *fileSystem.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment and one zone/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
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
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.FileStorageClient.GetFileSystem(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.FileSystem, nil
}

func getFileStorageFileSystemExports(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	var id string
	if h.Item != nil {
		id = getFileSystemID(h.Item)
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the API call to only the root compartment and one zone/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
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
		logger.Error("oci_file_storage_file_system.getFileStorageFileSystemExports", "connection_error", err)
		return nil, err
	}

	request := filestorage.ListExportsRequest{
		FileSystemId:  types.String(id),
		CompartmentId: &compartment,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.FileStorageClient.ListExports(ctx, request)
	if err != nil {
		logger.Error("oci_file_storage_file_system.getFileStorageFileSystemExports", "api_error", err)
		return nil, err
	}

	return response.Items, nil
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
	switch item := item.(type) {
	case filestorage.FileSystem:
		return item.FreeformTags
	case filestorage.FileSystemSummary:
		return item.FreeformTags
	}
	return nil
}

func fileSystemDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case filestorage.FileSystem:
		return item.DefinedTags
	case filestorage.FileSystemSummary:
		return item.DefinedTags
	}
	return nil
}

// Build additional filters
func buildFileStorageFileSystemFilters(equalQuals plugin.KeyColumnEqualsQualMap) filestorage.ListFileSystemsRequest {
	request := filestorage.ListFileSystemsRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}
	if equalQuals["id"] != nil {
		request.Id = types.String(equalQuals["id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = filestorage.ListFileSystemsLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}

func getFileSystemID(item interface{}) string {
	switch item := item.(type) {
	case filestorage.FileSystemSummary:
		return *item.Id
	case filestorage.FileSystem:
		return *item.Id
	}

	return ""
}