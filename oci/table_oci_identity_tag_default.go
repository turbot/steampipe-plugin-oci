package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableIdentityTagDefault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_tag_default",
		Description: "OCI Identity Tag Default",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIdentityTagDefault,
		},
		List: &plugin.ListConfig{
			Hydrate: listIdentityTagDefaults,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the tag default.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "tag_definition_name",
				Description: "The name used in the tag definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "tag_definition_id",
				Description: "The OCID of the tag definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "tag_namespace_id",
				Description: "The OCID of the tag namespace that contains the tag definition.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The tag default's current state.",
				Type:        proto.ColumnType_STRING,
			},

			// Other columns
			{
				Name:        "is_required",
				Description: "If you specify that a value is required, a value is set during resource creation (either by the user creating the resource or another tag default). If no value is set, resource creation is blocked.If the `isRequired` flag is set to true, the value is set during resource creation.If the `isRequired` flag is set to false, the value you enter is set during resource creation.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "time_created",
				Description: "Date and time the tagDefault was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "value",
				Description: "The default value for the tag definition.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},

			// Standard OCI columns
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listIdentityTagDefaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Trace("oci.listIdentityTagDefaults", "Compartment", compartment)

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.ListTagDefaultsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListTagDefaults(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, tagDefault := range response.Items {
			d.StreamListItem(ctx, tagDefault)
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

func getIdentityTagDefault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getIdentityTagDefault", "Compartment", compartment)

	// Rstrict the api call to only root compartment
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.GetTagDefaultRequest{
		TagDefaultId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetTagDefault(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.TagDefault, nil
}
