package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/objectstorage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableObjectStorageBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_objectstorage_bucket",
		Description: "OCI ObjectStorage Bucket",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"400", "404"}),
			Hydrate:    getObjectStorageBucket,
		},
		List: &plugin.ListConfig{
			Hydrate: listObjectStorageBuckets,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the bucket.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "namespace",
				Description: "The Object Storage namespace in which the bucket lives.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "approximate_count",
				Description: "The approximate number of objects in the bucket.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "approximate_size",
				Description: "The approximate total size in bytes of all objects in the bucket..",
				Type:        proto.ColumnType_INT,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "created_by",
				Description: "The OCID of the user who created the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The entity tag (ETag) for the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_read_only",
				Description: "Whether or not this bucket is read only.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of a master encryption key used to call the Key Management service to generate a data encryption key or to encrypt or decrypt a data encryption key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectStorageBucket,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "object_events_enabled",
				Description: "Whether or not events are emitted for object state changes in this bucket.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "object_lifecycle_policy_etag",
				Description: "The entity tag (ETag) for the live object lifecycle policy on the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "public_access_type",
				Description: "The type of public access enabled on this bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "replication_enabled",
				Description: "Whether or not this bucket is a replication source.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "storage_tier",
				Description: "The storage tier type assigned to the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "time_created",
				Description: "The date and time the bucket was created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "versioning",
				Description: "The versioning status on the bucket.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getObjectStorageBucket,
			},

			// json fields
			{
				Name:        "metadata",
				Description: "Arbitrary string keys and values for user-defined metadata.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectStorageBucket,
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectStorageBucket,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bucketTags),
				Hydrate:     getObjectStorageBucket,
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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

//// listObjectStorageBuckets FUNCTION
func listObjectStorageBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("listObjectStorageBuckets", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	nameSpace, err := getNamespace(ctx, d)
	if err != nil {
		return nil, err
	}

	request := objectstorage.ListBucketsRequest{
		CompartmentId: types.String(compartment),
		NamespaceName: &nameSpace.Value,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ObjectStorageClient.ListBuckets(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, bucketSummary := range response.Items {
			d.StreamListItem(ctx, bucketSummary)
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

func getObjectStorageBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getObjectStorageBucket")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("getObjectStorageBucket", "Compartment", compartment, "OCI_REGION", region)

	// Rstrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	var bucketName string
	if h.Item != nil {
		bucketName = *h.Item.(objectstorage.BucketSummary).Name
	} else {
		bucketName = d.KeyColumnQuals["name"].GetStringValue()
	}

	nameSpace, err := getNamespace(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := objectstorage.GetBucketRequest{
		NamespaceName: &nameSpace.Value,
		BucketName:    &bucketName,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}
	response, err := session.ObjectStorageClient.GetBucket(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Bucket, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func bucketTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var bucket objectstorage.Bucket
	var bucketSummary objectstorage.BucketSummary

	switch d.HydrateItem.(type) {
	case objectstorage.Bucket:
		bucket = d.HydrateItem.(objectstorage.Bucket)
	case objectstorage.BucketSummary:
		bucketSummary = d.HydrateItem.(objectstorage.BucketSummary)
	}

	var freeTags map[string]string
	var definedtags map[string]map[string]interface{}

	if bucket.Name != nil {
		freeTags = bucket.FreeformTags
		definedtags = bucket.DefinedTags
	} else {
		freeTags = bucketSummary.FreeformTags
		definedtags = bucketSummary.DefinedTags
	}

	var tags map[string]interface{}

	if freeTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeTags {
			tags[k] = v
		}
	}

	if definedtags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range definedtags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}

