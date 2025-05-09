package oci

import (
	"context"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type nameSpace struct {
	Value string
}

// // LIST FUNCTION
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
	request := objectstorage.GetNamespaceRequest{
		CompartmentId: &session.TenancyID,
	}

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
func getDefaultRetryPolicy(connection *plugin.Connection) *oci_common.RetryPolicy {
	// how many times to do the retry
	attempts := uint(9)
	minRetryDelay := 25 * time.Millisecond

	// Get config details for maximum error attempt and minimum delay time
	config := GetConfig(connection)
	if config.MaxErrorRetryAttempts != nil {
		attempts = uint(*config.MaxErrorRetryAttempts)
	}

	if config.MinErrorRetryDelay != nil {
		minRetryDelay = time.Duration(*config.MinErrorRetryDelay) * time.Millisecond
	}

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
			return (r.Error != nil && slices.Contains([]string{"429", "500", "503"}, statusCode))
		}
		return false
	}
	return getExponentialBackoffRetryPolicy(attempts, minRetryDelay, retryOnResponseCodes)
}

func getExponentialBackoffRetryPolicy(n uint, minRetryDelay time.Duration, fn func(r oci_common.OCIOperationResponse) bool) *oci_common.RetryPolicy {
	// the duration between each retry operation, you might want to waite longer each time the retry fails
	exponentialBackoff := func(r oci_common.OCIOperationResponse) time.Duration {

		// If errors are caused by load, retries can be ineffective if all API request retry at the same time.
		// To avoid this problem added a jitter of "+/-20%" with delay time.
		// For example, if the delay is 25ms, the final delay could be between 20 and 30ms.
		var jitter = float64(rand.Intn(120-80)+80) / 100

		// Creates a new exponential backoff using the starting value of
		// minDelay and (minDelay * 3^retrycount) * jitter on each failure
		// as example (23.25ms, 63ms, 238.5ms, 607.4ms, 2s, 5.22s, 20.31s...) up to max.
		// Maximum delay should not be more than 3 min
		maxDelayTime := time.Duration(int(float64(int(minRetryDelay.Nanoseconds())*int(math.Pow(3, float64(r.AttemptNumber)))) * jitter))
		if maxDelayTime > time.Duration(3*time.Minute) {
			return time.Duration(3 * time.Minute)
		}

		return maxDelayTime
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

func extractTags(freeformTags map[string]string, definedTags map[string]map[string]interface{}) map[string]interface{} {
	tags := make(map[string]interface{})

	for k, v := range freeformTags {
		tags[k] = v
	}

	for _, v := range definedTags {
		for key, value := range v {
			tags[key] = value
		}
	}

	return tags
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getTenantIdMemoized = plugin.HydrateFunc(getTenantIdUncached).Memoize(memoize.WithCacheKeyFunction(getTenantIdCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getTenantId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getTenantIdMemoized(ctx, d, h)
}

// Build a cache key for the call to getTenantId, including the region since this is a multi-region call.
// Notably, this may be called WITHOUT a region. In that case we just share a cache for non-region data.
func getTenantIdCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getTenantId"
	return key, nil
}

func getTenantIdUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getTenantId")

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	return session.TenancyID, nil
}
