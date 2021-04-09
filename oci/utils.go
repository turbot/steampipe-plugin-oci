package oci

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/objectstorage"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type nameSpace struct {
	Value string
}

//// LIST FUNCTION
func getNamespace(ctx context.Context, d *plugin.QueryData, region string) (*nameSpace, error) {
	plugin.Logger(ctx).Trace("getNamespace")

	cacheKey := "ObjectStorageNamespace"

	// check if the namespace is already saved in cache
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*nameSpace), nil
	}

	// Create Session
	session, err := objectStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	request := objectstorage.GetNamespaceRequest{}

	response, err := session.ObjectStorageClient.GetNamespace(ctx, request)
	if err != nil {
		return nil, err
	}
	name := &nameSpace{
		Value: *response.Value,
	}
	d.ConnectionManager.Cache.Set(cacheKey, name)

	return name, err
}

// https://github.com/oracle/oci-go-sdk/blob/master/example/helpers/helper.go#L127
func getDefaultRetryPolicy() *oci_common.RetryPolicy {
	// how many times to do the retry
	attempts := uint(5)

	/*
		429	TooManyRequests	You have issued too many requests to the
		Oracle Cloud Infrastructure APIs in too short of an amount of time.	Yes, with backoff.

		500	InternalServerError	An internal server error occurred.	Yes, with backoff.

		503	ServiceUnavailable	The service is currently unavailable.	Yes, with backoff.
		https: //docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm
	*/
	retryOnResponseCodes := func(r oci_common.OCIOperationResponse) bool {
		if r.Response != nil && r.Response.HTTPResponse() != nil {
			statusCode := strconv.Itoa(r.Response.HTTPResponse().StatusCode)
			return (r.Error != nil && helpers.StringSliceContains([]string{"429", "500", "503"}, statusCode))
		}
		return false
	}
	return getExponentialBackoffRetryPolicy(attempts, retryOnResponseCodes)
}

func getExponentialBackoffRetryPolicy(n uint, fn func(r oci_common.OCIOperationResponse) bool) *oci_common.RetryPolicy {
	// the duration between each retry operation, you might want to waite longer each time the retry fails
	exponentialBackoff := func(r oci_common.OCIOperationResponse) time.Duration {
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}
	policy := oci_common.NewRetryPolicy(n, fn, exponentialBackoff)
	return &policy
}

// Extract oci region name from resource oci id
func ociRegionName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	splittedID := strings.Split(id, ".")
	return oci_common.StringToRegion(types.SafeString(splittedID[3])), nil
}
