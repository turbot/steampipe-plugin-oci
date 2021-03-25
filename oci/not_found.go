package oci

import (
	"strconv"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// function which returns an ErrorPredicate for AWS API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		if ociErr, ok := err.(oci_common.ServiceError); ok {
			return helpers.StringSliceContains(notFoundErrors, strconv.Itoa(ociErr.GetHTTPStatusCode()))
		}
		return false
	}
}
