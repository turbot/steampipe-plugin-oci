package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/objectstorage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableObjectStorageObject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_objectstorage_object",
		Description: "OCI Object Storage Object",
		// Object can have same name in two different buckets, regions or compartments, leading to duplicate result in get call
		List: &plugin.ListConfig{
			Hydrate:       listObjectStorageObjects,
			ParentHydrate: listObjectStorageBuckets,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bucket_name",
				Description: "The name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "namespace",
				Description: "The Object Storage namespace used for the request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "archival_state",
				Description: "Archival state of an object. This field is set only for objects in Archive tier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cache_control",
				Description: "The Cache-Control header.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_disposition",
				Description: "The Content-Disposition header.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_encoding",
				Description: "The Content-Encoding header.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_language",
				Description: "The Content-Language header.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_range",
				Description: "Content-Range header for range requests.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_type",
				Description: "The Content-Type header.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The current entity tag (ETag) for the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expires",
				Description: "The date and time after which the object is no longer cached by a browser, proxy, or other caching entity.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Expires.Time"),
			},
			{
				Name:        "is_not_modified",
				Description: "Flag to indicate whether or not the object was modified. If this is true, the getter for the object itself will return null.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_BOOL,
				Default:     false,
			},
			{
				Name:        "md5",
				Description: "Base64-encoded MD5 hash of the object data.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "opc_multipart_md5",
				Description: "Base-64 representation of the multipart object hash. Only applicable to objects uploaded using multipart upload.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size",
				Description: "Size of the object in bytes.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "storage_tier",
				Description: "The storage tier that the object is stored in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the object was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_modified",
				Description: "The date and time the object was modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeModified.Time"),
			},
			{
				Name:        "time_of_archival",
				Description: "Time that the object is returned to the archived state. This field is only present for restored objects.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfArchival.Time"),
			},
			{
				Name:        "version_id",
				Description: "The version ID of the object requested.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionId"),
			},
			{
				Name:        "opc_meta",
				Description: "The user-defined metadata for the object.",
				Hydrate:     getObjectStorageObject,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
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

type objectInfo struct {
	BucketName string
	Namespace  string
	Region     string
	objectstorage.ObjectSummary
}

//// LIST FUNCTION

func listObjectStorageObjects(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Trace("listObjectStorageObjects", "OCI_REGION", region)

	bucketName := *h.Item.(bucketInfo).Name

	objectNameSpace, err := getNamespace(ctx, d, region)
	if err != nil {
		logger.Error("listObjectStorageObjects", "error_getNamespace", region)
		return nil, err
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := objectstorage.ListObjectsRequest{
		BucketName:    &bucketName,
		NamespaceName: &objectNameSpace.Value,
		Fields:        types.String("name,size,etag,timeCreated,md5,timeModified,storageTier,archivalState"),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	response, err := session.ObjectStorageClient.ListObjects(ctx, request)
	if err != nil {
		logger.Error("listObjectStorageObjects", "error_ListObjects", err)
		return nil, err
	}

	for _, objectSummary := range response.Objects {
		d.StreamListItem(ctx, objectInfo{bucketName, objectNameSpace.Value, region, objectSummary})
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getObjectStorageObject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getObjectStorageObject")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getObjectStorageObject", "OCI_REGION", region)

	var bucketName, namespace, objectName string
	if h.Item != nil {
		info := h.Item.(objectInfo)
		bucketName = info.BucketName
		objectName = *info.Name
		namespace = info.Namespace
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		logger.Error("getObjectStorageObject", "error_objectStorageService", err)
		return nil, err
	}

	request := objectstorage.GetObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ObjectStorageClient.GetObject(ctx, request)
	if err != nil {
		logger.Error("getObjectStorageObject", "error_GetObject", err)
		if ociErr, ok := err.(common.ServiceError); ok {
			if ociErr.GetCode() == "NotRestored" {
				return nil, nil
			}
		}
		return nil, err
	}

	return response, nil
}
