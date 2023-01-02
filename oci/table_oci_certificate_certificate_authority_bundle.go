package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/certificates"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesCertificateAuthorityBundle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_certificate_authority_bundle",
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
					Require: plugin.Optional,
				},
				{
					Name:    "version_name",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCertificatesCertificateAuthorityBundle,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_authority_id",
					Require: plugin.Required,
				},
				{
					Name:    "version_number",
					Require: plugin.Optional,
				},
				{
					Name:    "version_name",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCertificatesCertificateAuthorityBundle,
		},
		GetMatrixItemFunc: BuildRegionList,
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
			},
			{
				Name:        "version_number",
				Description: "The version number of the CA.",
				Type:        proto.ColumnType_INT,
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
				Name:        "cert_chain_pem",
				Description: "The certificate chain (in PEM format) for this CA version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_name",
				Description: "The name of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revocation_status",
				Description: "The revocation status for the certificate authority.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Authority Bundle was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

// HYDRATE FUNCTION
func getCertificatesCertificateAuthorityBundle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getCertificatesCertificateAuthorityBundle", "OCI_REGION", region)

	request := buildGetCertificatesCertificateAuthorityBundleFilters(d.KeyColumnQuals)

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesCertificateAuthorityBundle", "error_CertificatesService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesClient.GetCertificateAuthorityBundle(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.CertificateAuthorityBundle, nil
}

// Build additional filters
func buildGetCertificatesCertificateAuthorityBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificates.GetCertificateAuthorityBundleRequest {
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
