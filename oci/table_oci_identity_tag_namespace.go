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

func tableIdentityTagNamespace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_identity_tag_namespace",
		Description: "OCI Identity Tag Namespace",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getIdentityTagNamespace,
		},
		List: &plugin.ListConfig{
			Hydrate: listIdentityTagNamespaces,
		},
		GetMatrixItem: BuildCompartmentList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the tag namespace. It must be unique across all tag namespaces in the tenancy and cannot be changed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the tag namespace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "is_retired",
				Description: "Whether the tag namespace is retired.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "lifecycle_state",
				Description: "The tagnamespace's current state.",
				Type:        proto.ColumnType_STRING,
			},

			// Other columns
			{
				Name:        "time_created",
				Description: "Date and time the tagNamespace was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "description",
				Description: "The description you assign to the tag namespace.",
				Type:        proto.ColumnType_STRING,
			},

			// Tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(identityTagNamespaceTags),
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

func listIdentityTagNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Trace("oci.listIdentityTagNamespaces", "Compartment", compartment)

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.ListTagNamespacesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListTagNamespaces(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, tagNamespace := range response.Items {
			d.StreamListItem(ctx, tagNamespace)
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

func getIdentityTagNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getIdentityTagNamespace", "Compartment", compartment)

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

	request := identity.GetTagNamespaceRequest{
		TagNamespaceId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.IdentityClient.GetTagNamespace(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.TagNamespace, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 2. Defined Tags
// 3. Free-form tags
func identityTagNamespaceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case identity.TagNamespaceSummary:
		tagNamespace := d.HydrateItem.(identity.TagNamespaceSummary)
		freeformTags = tagNamespace.FreeformTags
		definedTags = tagNamespace.DefinedTags
	case identity.TagNamespace:
		tagNamespace := d.HydrateItem.(identity.TagNamespace)
		freeformTags = tagNamespace.FreeformTags
		definedTags = tagNamespace.DefinedTags
	}

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	if definedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range definedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
