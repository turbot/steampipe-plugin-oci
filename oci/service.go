package oci

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type session struct {
	TenancyID      string
	IdentityClient identity.IdentityClient
	ComputeClient  core.ComputeClient
}

// identityService returns the service client for OCI Identity service
func identityService(ctx context.Context, d *plugin.QueryData) (*session, error) {
	// if region == "" {
	// 	return nil, fmt.Errorf("region must be passed ACMService")
	// }
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("Identity-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider := oci_common.CustomProfileConfigProvider(*ociConfig.ConfigPath, *ociConfig.Profile)
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:      tenantId,
		IdentityClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// coreComputeService returns the service client for OCI Core Compute service
func coreComputeService(ctx context.Context, d *plugin.QueryData) (*session, error) {
	// if region == "" {
	// 	return nil, fmt.Errorf("region must be passed ACMService")
	// }
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("Compute-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider := oci_common.CustomProfileConfigProvider(*ociConfig.ConfigPath, *ociConfig.Profile)
	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:     tenantId,
		ComputeClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// https://github.com/oracle/oci-go-sdk/blob/master/example/helpers/helper.go#L127
func getDefaultRetryPolicy() *oci_common.RetryPolicy {
	// how many times to do the retry
	attempts := uint(5)

	// 429	TooManyRequests	You have issued too many requests to the
	// Oracle Cloud Infrastructure APIs in too short of an amount of time.	Yes, with backoff.

	// 500	InternalServerError	An internal server error occurred.	Yes, with backoff.

	// 503	ServiceUnavailable	The service is currently unavailable.	Yes, with backoff.
	// https: //docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm
	retryOnResponseCodes := func(r oci_common.OCIOperationResponse) bool {
		statusCode := strconv.Itoa(r.Response.HTTPResponse().StatusCode)
		return (r.Error != nil && helpers.StringSliceContains([]string{"429", "500", "503"}, statusCode))
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
