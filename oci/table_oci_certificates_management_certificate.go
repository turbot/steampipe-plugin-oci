package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesManagementCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificates_management_certificate",
		Description:      "OCI Management Certificate",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCertificatesManagementCertificate,
		},
		List: &plugin.ListConfig{
			Hydrate: listCertificatesManagementCertificates,
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
				Description: "The OCID of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "A user-friendly name for the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "config_type",
				Description: "The origin of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "issuer_certificate_authority_id",
				Description: "The OCID of the certificate authority (CA) that issued the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A brief description of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_rules",
				Description: "A list of rules that control how the certificate is used and managed.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the certificate version.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycle state of the certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCertificatesManagementCertificate,
			},
			{
				Name:        "current_version",
				Description: "Details about the current version of the certificate.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCertificatesManagementCertificate,
			},
			{
				Name:        "subject",
				Description: "Certificate subject informnation.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificate_revocation_list_details",
				Description: "Certificate revocation details.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCertificatesManagementCertificate,
			},
			{
				Name:        "key_algorithm",
				Description: "The algorithm used to create key pairs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signature_algorithm",
				Description: "The algorithm used to sign the public key certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_profile_type",
				Description: "The name of the profile used to create the certificate.",
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
				Description: "Time that the Certificate was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(certificatesManagementCertificateTags),
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCertificatesManagementCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_certificates_management_certificate.listCertificatesManagementCertificates", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("oci_certificates_management_certificate.listCertificatesManagementCertificates", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildListCertificatesManagementCertificateFilters(equalQuals)
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
		response, err := session.CertificatesManagementClient.ListCertificates(ctx, request)
		if err != nil {
			plugin.Logger(ctx).Error("oci_certificates_management_certificate.listCertificatesManagementCertificates", "api_error", err)
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

//// HYDRATE FUNCTIONS

func getCertificatesManagementCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_certificates_management_certificate.getCertificatesManagementCertificate", "Compartment", compartment, "OCI_REGION", region)
	if h.Item == nil && !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	request := buildGetCertificatesManagementCertificateFilters(d.EqualsQuals, h)

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		logger.Error("oci_certificates_management_certificate.getCertificatesManagementCertificate", "connection_error", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesManagementClient.GetCertificate(ctx, request)
	if err != nil {
		logger.Error("oci_certificates_management_certificate.getCertificatesManagementCertificate", "api_error", err)
		return nil, err
	}
	return response.Certificate, nil
}

//// TRANSFORM FUNCTION

func certificatesManagementCertificateTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case certificatesmanagement.Certificate:
		obj := d.HydrateItem.(certificatesmanagement.Certificate)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case certificatesmanagement.CertificateSummary:
		obj := d.HydrateItem.(certificatesmanagement.CertificateSummary)
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
func buildListCertificatesManagementCertificateFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificatesmanagement.ListCertificatesRequest {
	request := certificatesmanagement.ListCertificatesRequest{}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = certificatesmanagement.ListCertificatesLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
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
func buildGetCertificatesManagementCertificateFilters(equalQuals plugin.KeyColumnEqualsQualMap, h *plugin.HydrateData) certificatesmanagement.GetCertificateRequest {
	request := certificatesmanagement.GetCertificateRequest{}

	if h.Item != nil {
		request.CertificateId = h.Item.(certificatesmanagement.CertificateSummary).Id
	} else {
		request.CertificateId = types.String(equalQuals["id"].GetStringValue())
	}

	return request
}
