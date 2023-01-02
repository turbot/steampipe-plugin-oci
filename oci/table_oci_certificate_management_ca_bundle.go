package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strings"
)

// TABLE DEFINITION
func tableCertificatesManagementCaBundle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_management_ca_bundle",
		Description:      "OCI Ca Bundle",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCertificatesManagementCaBundle,
		},
		List: &plugin.ListConfig{
			Hydrate: listCertificatesManagementCaBundles,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the CA bundle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "A user-friendly name for the CA bundle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the CA bundle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A brief description of the CA bundle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycle state of the CA bundle.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "freeform_tags",
				Description: "Free-form tags for this resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_created",
				Description: "Time that the Ca Bundle was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(certificatesManagementCaBundleTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

// LIST FUNCTION
func listCertificatesManagementCaBundles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCertificatesManagementCaBundles", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildListCertificatesManagementCaBundleFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(20)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.CertificatesManagementClient.ListCaBundles(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

// HYDRATE FUNCTION
func getCertificatesManagementCaBundle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCertificatesManagementCaBundle", "Compartment", compartment, "OCI_REGION", region)
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	request := buildGetCertificatesManagementCaBundleFilters(d.KeyColumnQuals, h)

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesManagementCaBundle", "error_CertificatesManagementService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesManagementClient.GetCaBundle(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.CaBundle, nil
}

// TRANSFORM FUNCTION
func certificatesManagementCaBundleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case certificatesmanagement.CaBundle:
		obj := d.HydrateItem.(certificatesmanagement.CaBundle)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case certificatesmanagement.CaBundleSummary:
		obj := d.HydrateItem.(certificatesmanagement.CaBundleSummary)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
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

// Build additional list filters
func buildListCertificatesManagementCaBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificatesmanagement.ListCaBundlesRequest {
	request := certificatesmanagement.ListCaBundlesRequest{}

	if equalQuals["compartment_id"] != nil {
		request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = certificatesmanagement.ListCaBundlesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = certificatesmanagement.ListCaBundlesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}

	return request
}

// Build additional filters
func buildGetCertificatesManagementCaBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap, h *plugin.HydrateData) certificatesmanagement.GetCaBundleRequest {
	request := certificatesmanagement.GetCaBundleRequest{}

	if h.Item != nil {
		request.CaBundleId = h.Item.(certificatesmanagement.CaBundleSummary).Id
	} else {
		request.CaBundleId = types.String(equalQuals["id"].GetStringValue())
	}

	return request
}
