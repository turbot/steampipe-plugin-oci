package oci

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/objectstorage"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
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
	attempts := uint(9)

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
		// Minumum delay time
		var minRetryDelay time.Duration = 25 * time.Millisecond

		// If errors are caused by load, retries can be ineffective if all API request retry at the same time.
		// To avoid this problem added a jitter of "+/-20%" with delay time.
		// For example, if the delay is 25ms, the final delay could be between 20 and 30ms.
		var jitter = float64(rand.Intn(120-80)+80) / 100

		// Creates a new exponential backoff using the starting value of
		// minDelay and (minDelay * 3^retrycount) * jitter on each failure
		// as example (23.25ms, 63ms, 238.5ms, 607.4ms, 2s, 5.22s, 20.31s...) up to max.
		return time.Duration(int(float64(int(minRetryDelay.Nanoseconds())*int(math.Pow(3, float64(r.AttemptNumber)))) * jitter))
	}
	policy := oci_common.NewRetryPolicy(n, fn, exponentialBackoff)
	return &policy
}

// Extract OCI region name from the resource id
func ociRegionName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	splittedID := strings.Split(id, ".")
	regionName := oci_common.StringToRegion(types.SafeString(splittedID[3]))
	return regionName, nil
}
