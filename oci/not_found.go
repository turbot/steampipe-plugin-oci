package oci

import (
	"context"
	"strconv"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// function which returns an ErrorPredicate for OCI API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if ociErr, ok := err.(oci_common.ServiceError); ok {
			return helpers.StringSliceContains(notFoundErrors, strconv.Itoa(ociErr.GetHTTPStatusCode()))
		}
		return false
	}
}

// function which returns an ErrorPredicateWithContext for OCI API calls
// https://docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm
// It's advisable to handle errors based on their error codes rather than relying solely on the HTTP status code. This is because different errors can have the same HTTP status code, but they will have distinct error codes.
func isNotFoundErrorCode(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if ociErr, ok := err.(oci_common.ServiceError); ok {
			return helpers.StringSliceContains(notFoundErrors, ociErr.GetCode())
		}
		return false
	}
}
