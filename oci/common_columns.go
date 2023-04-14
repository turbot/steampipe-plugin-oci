package oci

import (
	"context"

	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const (
	// Constants for Standard Column Descriptions
	ColumnDescriptionAkas        = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags        = "A map of tags for the resource."
	ColumnDescriptionTitle       = "Title of the resource."
	ColumnDescriptionTenantId    = "The OCID of the Tenant in which the resource is located."
	ColumnDescriptionTenantName  = "The name of the Tenant in which the resource is located."
	ColumnDescriptionCompartment = "The OCID of the compartment in Tenant in which the resource is located."
	ColumnDescriptionRegion      = "The OCI region in which the resource is located."

	// Other repetitive columns for the provider
	ColumnDescriptionFreefromTags = "Free-form tags for resource. This tags can be applied by any user with permissions on the resource."
	ColumnDescriptionDefinedTags  = "Defined tags for resource. Defined tags are set up in your tenancy by an administrator. Only users granted permission to work with the defined tags can apply them to resources."
	ColumnDescriptionSystemTags   = "System tags for resource. System tags can be viewed by users, but can only be created by the system."
)

func commonColumnsForAllResource(columns []*plugin.Column) []*plugin.Column {
	return append(columns, []*plugin.Column{
		{
			Name:        "tenant_name",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCommonColumns,
			Description: ColumnDescriptionTenantName,
			Transform:   transform.FromField("Name"),
		},
	}...)
}

// returns the tenant_name common column which is added across all the tables
func getCommonColumns(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Trace logging to debug cache and execution flows
	plugin.Logger(ctx).Debug("getCommonColumns", "status", "starting", "connection_name", d.Connection.Name)

	var tenancyData identity.Tenancy

	tenancy, err := getTenancyData(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("getCommonColumns", "status", "failed", "connection_name", d.Connection.Name, "error", err)
		return nil, err
	}

	tenancyData = tenancy.(identity.Tenancy)
        
	plugin.Logger(ctx).Debug("getCommonColumns", "status", "finished", "connection_name", d.Connection.Name, "tenancyData", tenancyData)

	return tenancyData, nil
}

var getTenancyData = plugin.HydrateFunc(getTenancydataUncached).Memoize()

func getTenancydataUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getTenancydataUncached", "status", "failed", "connection_name", d.Connection.Name, "session_error", err)
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.GetTenancyRequest{
		TenancyId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.GetTenancy(ctx, request)
	if err != nil {
		plugin.Logger(ctx).Error("getTenancydataUncached", "status", "failed", "connection_name", d.Connection.Name, "api_error", err)
		return nil, err
	}

	return response.Tenancy, nil
}
