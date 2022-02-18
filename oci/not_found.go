package oci

import (
	"strconv"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
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
