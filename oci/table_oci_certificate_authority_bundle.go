package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/certificates"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesAuthorityBundle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_authority_bundle",
		Description:      "OCI Certificate Authority Bundle",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_authority_id",
					Require: plugin.Required,
				},
				{
					Name:    "version_number",
					Require: plugin.Required,
				},
				{
					Name:    "version_name",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCertificateAuthorityBundle,
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "certificate_authority_id",
				Description: "The OCID of the certificate authority (CA).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_authority_name",
				Description: "The name of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "serial_number",
				Description: "A unique certificate identifier used in certificate revocation tracking, formatted as octets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_pem",
				Description: "The certificate (in PEM format) for this CA version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCertificateAuthorityBundle,
			},
			{
				Name:        "version_number",
				Description: "The version number of the CA.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cert_chain_pem",
				Description: "The certificate chain (in PEM format) for this CA version.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCertificateAuthorityBundle,
			},
			{
				Name:        "version_name",
				Description: "The name of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Authority Bundle was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "validity",
				Description: "Validatity details for the certificate authority.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "stages",
				Description: "A list of rotation states for this CA.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "revocation_status",
				Description: "The revocation status for the certificate authority.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateAuthorityName"),
			},
			// Standard OCI columns
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

//// HYDRATE FUNCTIONS

func getCertificateAuthorityBundle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_certificate_authority_bundle.getCertificateAuthorityBundle", "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	request := buildGetCertificateAuthorityBundleFilters(d.EqualsQuals)

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_certificate_authority_bundle.getCertificateAuthorityBundle", "connection_error", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesClient.GetCertificateAuthorityBundle(ctx, request)
	if err != nil {
		logger.Error("oci_certificate_authority_bundle.getCertificateAuthorityBundle", "api_error", err)
		return nil, err
	}
	return response.CertificateAuthorityBundle, nil
}

// Build additional filters
func buildGetCertificateAuthorityBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificates.GetCertificateAuthorityBundleRequest {
	request := certificates.GetCertificateAuthorityBundleRequest{}

	if equalQuals["certificate_authority_id"] != nil && equalQuals["certificate_authority_id"].GetStringValue() != "" {
		request.CertificateAuthorityId = types.String(equalQuals["certificate_authority_id"].GetStringValue())
	}
	if equalQuals["version_number"] != nil && equalQuals["version_number"].GetInt64Value() != 0 {
		request.VersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}
	if equalQuals["version_name"] != nil && equalQuals["version_name"].GetStringValue() != "" {
		request.CertificateAuthorityVersionName = types.String(equalQuals["version_name"].GetStringValue())
	}

	return request
}
