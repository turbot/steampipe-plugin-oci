package oci

import (
	"context"
	"strconv"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableObjectStorageBucket(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_objectstorage_bucket",
		Description: "OCI ObjectStorage Bucket",
		// Bucket can have same name in two different compartments, leads to duplicate result in get call
		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.SingleColumn("name"),
		// 	ShouldIgnoreError: isNotFoundError([]string{"400", "404"}),
		// 	Hydrate:           getObjectStorageBucket,
		// },
		List: &plugin.ListConfig{
			Hydrate: listObjectStorageBuckets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "namespace",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
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
				Description: "The approximate total size in bytes of all objects in the bucket.",
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
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
			{
				Name:        "object_lifecycle_policy",
				Description: "Specifies the object lifecycle policy for the bucket.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getObjectStorageBucketObjectLifecycle,
				Transform:   transform.FromValue(),
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
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
			},
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

type bucketInfo struct {
	Region string
	objectstorage.BucketSummary
}

//// LIST FUNCTION

func listObjectStorageBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Error("listObjectStorageBuckets", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	nameSpace, err := getNamespace(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := objectstorage.ListBucketsRequest{
		CompartmentId: types.String(compartment),
		NamespaceName: &nameSpace.Value,
		Limit:         types.Int(1000),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	if equalQuals["namespace"] != nil {
		request.NamespaceName = types.String(equalQuals["namespace"].GetStringValue())
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ObjectStorageClient.ListBuckets(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, bucketSummary := range response.Items {
			d.StreamListItem(ctx, bucketInfo{region, bucketSummary})

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

func getObjectStorageBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getObjectStorageBucket")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getObjectStorageBucket", "Compartment", compartment, "OCI_REGION", region)

	var bucketName, nameSpace string
	if h.Item != nil {
		bucket := h.Item.(bucketInfo)
		bucketName = *bucket.Name
		nameSpace = *bucket.Namespace
	} else {
		bucketName = d.EqualsQuals["name"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty bucket name in get call
	if bucketName == "" {
		return nil, nil
	}

	if nameSpace == "" {
		bucketNamespace, err := getNamespace(ctx, d, region)
		if err != nil {
			return nil, err
		}
		nameSpace = bucketNamespace.Value
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := objectstorage.GetBucketRequest{
		NamespaceName: &nameSpace,
		BucketName:    &bucketName,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.ObjectStorageClient.GetBucket(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Bucket, nil
}

func getObjectStorageBucketObjectLifecycle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getObjectStorageBucketObjectLifecycle")

	data := h.Item.(bucketInfo)
	bucketName := data.Name
	nameSpace := data.Namespace

	// Create Session
	session, err := objectStorageService(ctx, d, data.Region)
	if err != nil {
		return nil, err
	}

	request := objectstorage.GetObjectLifecyclePolicyRequest{
		NamespaceName: nameSpace,
		BucketName:    bucketName,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	response, err := session.ObjectStorageClient.GetObjectLifecyclePolicy(ctx, request)
	if err != nil {
		if ociErr, ok := err.(oci_common.ServiceError); ok {
			if helpers.StringSliceContains([]string{"LifecyclePolicyNotFound"}, strconv.Itoa(ociErr.GetHTTPStatusCode())) {
				return nil, nil
			}
		}
	}

	return response.ObjectLifecyclePolicy, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func bucketTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	bucket := d.HydrateItem.(objectstorage.Bucket)

	var tags map[string]interface{}
	if bucket.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range bucket.FreeformTags {
			tags[k] = v
		}
	}

	if bucket.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range bucket.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
