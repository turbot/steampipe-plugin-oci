package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v36/objectstorage"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const (
	// Constants for Standard Column Descriptions
	ColumnDescriptionAkas        = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags        = "A map of tags for the resource."
	ColumnDescriptionTitle       = "Title of the resource."
	ColumnDescriptionTenant      = "The OCID of the Tenant in which the resource is located."
	ColumnDescriptionCompartment = "The OCID of the compartment in Tenant in which the resource is located."
	ColumnDescriptionRegion      = "The OCI region in which the resource is located."

	// Other repetitive columns for the provider
	ColumnDescriptionFreefromTags = "Free-form tags for resource. This tags can be applied by any user with permissions on the resource."
	ColumnDescriptionDefinedTags  = "Defined tags for resource. Defined tags are set up in your tenancy by an administrator. Only users granted permission to work with the defined tags can apply them to resources."
)

type nameSpace struct {
	Value string
}

//// listObjectStorageBuckets FUNCTION
func getNamespace(ctx context.Context, d *plugin.QueryData, region string) (*nameSpace, error) {
	plugin.Logger(ctx).Trace("getNamespace")
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("getNamespace", "Compartment", compartment, "OCI_REGION", region)

	serviceCacheKey := fmt.Sprintf("Namespace-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
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
	d.ConnectionManager.Cache.Set(serviceCacheKey, name)

	return name, err
}
