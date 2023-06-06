package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

// TABLE DEFINITION
func tableCertificatesManagementCertificateAuthority(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_management_certificate_authority",
		Description:      "OCI Certificate Authority",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCertificatesManagementCertificateAuthority,
		},
		List: &plugin.ListConfig{
			Hydrate: listCertificatesManagementCertificateAuthorities,
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
				{
					Name:    "issuer_certificate_authority_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "A user-friendly name for the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Authority was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the certificate authority.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "config_type",
				Description: "The origin of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "issuer_certificate_authority_id",
				Description: "The OCID of the parent CA that issued this CA. If this is the root CA, then this value is null.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A brief description of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the CA version.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the Oracle Cloud Infrastructure Vault key used to encrypt the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current CA lifecycle state.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCertificatesManagementCertificateAuthority,
			},
			{
				Name:        "signing_algorithm",
				Description: "The algorithm used to sign public key certificates that the CA issues.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_authority_rules",
				Description: "An optional list of rules that control how the CA is used and managed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "current_version",
				Description: "Details about the current version of the CA.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCertificatesManagementCertificateAuthority,
			},
			{
				Name:        "certificate_revocation_list_details",
				Description: "Details about the CA revocation list.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCertificatesManagementCertificateAuthority,
			},
			{
				Name:        "subject",
				Description: "Subject information for the CA.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: "Simple key-value pair that is applied without any predefined name.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Usage of predefined tag keys.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(certificatesManagementCertificateAuthorityTags),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCertificatesManagementCertificateAuthorities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listCertificatesManagementCertificateAuthorities", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
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
	request := buildListCertificatesManagementCertificateAuthorityFilters(equalQuals)
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
		response, err := session.CertificatesManagementClient.ListCertificateAuthorities(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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

//// HYDRATE FUNCTION

func getCertificatesManagementCertificateAuthority(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getCertificatesManagementCertificateAuthority", "Compartment", compartment, "OCI_REGION", region)
	if h.Item == nil && !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	request := buildGetCertificatesManagementCertificateAuthorityFilters(d.EqualsQuals, h)

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesManagementCertificateAuthority", "error_CertificatesManagementService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesManagementClient.GetCertificateAuthority(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.CertificateAuthority, nil
}

//// TRANSFORM FUNCTION

func certificatesManagementCertificateAuthorityTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case certificatesmanagement.CertificateAuthority:
		obj := d.HydrateItem.(certificatesmanagement.CertificateAuthority)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case certificatesmanagement.CertificateAuthoritySummary:
		obj := d.HydrateItem.(certificatesmanagement.CertificateAuthoritySummary)
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
func buildListCertificatesManagementCertificateAuthorityFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificatesmanagement.ListCertificateAuthoritiesRequest {
	request := certificatesmanagement.ListCertificateAuthoritiesRequest{}

	if equalQuals["compartment_id"] != nil {
		request.CompartmentId = types.String(equalQuals["compartment_id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = certificatesmanagement.ListCertificateAuthoritiesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = certificatesmanagement.ListCertificateAuthoritiesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}

	if equalQuals["issuer_certificate_authority_id"] != nil {
		request.IssuerCertificateAuthorityId = types.String(equalQuals["issuer_certificate_authority_id"].GetStringValue())
	}

	return request
}

// Build additional filters
func buildGetCertificatesManagementCertificateAuthorityFilters(equalQuals plugin.KeyColumnEqualsQualMap, h *plugin.HydrateData) certificatesmanagement.GetCertificateAuthorityRequest {
	request := certificatesmanagement.GetCertificateAuthorityRequest{}

	if h.Item != nil {
		request.CertificateAuthorityId = h.Item.(certificatesmanagement.CertificateAuthoritySummary).Id
	} else {
		request.CertificateAuthorityId = types.String(equalQuals["id"].GetStringValue())
	}

	return request
}