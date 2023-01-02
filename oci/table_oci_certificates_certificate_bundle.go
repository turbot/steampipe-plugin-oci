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
func tableCertificatesCertificateBundle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificates_certificate_bundle",
		Description:      "OCI Certificate Bundle",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_id",
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
				{
					Name:    "certificate_bundle_type",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCertificatesCertificateBundle,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_id",
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
				{
					Name:    "certificate_bundle_type",
					Require: plugin.Optional,
				},
			},
			Hydrate: getCertificatesCertificateBundle,
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "certificate_id",
				Description: "The certificate OCID.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_name",
				Description: "The friendly name of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_number",
				Description: "The certificate version number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "serial_number",
				Description: "The certificate serial number.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "validity",
				Description: "The validity details about the certificate.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "stages",
				Description: "The certificate stages.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificate_pem",
				Description: "The public certificate in pem format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cert_chain_pem",
				Description: "The certificate chain in pem format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_key_pem",
				Description: "The private key in pem format.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_name",
				Description: "A friendly name for the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "revocation_status",
				Description: "Revocation status of the certificate.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificate_bundle_type",
				Description: "The type of certificate bundle.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("certificate_bundle_type"),
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Bundle was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CertificateName"),
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
func getCertificatesCertificateBundle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getCertificatesCertificateBundle", "OCI_REGION", region)

	request := buildGetCertificatesCertificateBundleFilters(d.KeyColumnQuals)

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesCertificateBundle", "error_CertificatesService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesClient.GetCertificateBundle(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.CertificateBundle, nil
}

// Build additional filters
func buildGetCertificatesCertificateBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificates.GetCertificateBundleRequest {
	request := certificates.GetCertificateBundleRequest{}

	if equalQuals["certificate_id"] != nil && equalQuals["certificate_id"].GetStringValue() != "" {
		request.CertificateId = types.String(equalQuals["certificate_id"].GetStringValue())
	}
	if equalQuals["version_number"] != nil && equalQuals["version_number"].GetInt64Value() != 0 {
		request.VersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}
	if equalQuals["version_name"] != nil && equalQuals["version_name"].GetStringValue() != "" {
		request.CertificateVersionName = types.String(equalQuals["version_name"].GetStringValue())
	}
	if equalQuals["certificate_bundle_type"] != nil && equalQuals["certificate_bundle_type"].GetStringValue() != "" {
		request.CertificateBundleType = certificates.GetCertificateBundleCertificateBundleTypeEnum(equalQuals["certificate_bundle_type"].GetStringValue())
	}

	return request
}
