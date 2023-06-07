package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/certificates"
	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesCertificateBundleVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_certificate_bundle_version",
		Description:      "OCI Certificate Bundle Version",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_id",
					Require: plugin.Required,
				},
				{
					Name:    "version_number",
					Require: plugin.Required,
				},
			},
			Hydrate: getCertificatesCertificateBundleVersion,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCertificatesManagementCertificates,
			Hydrate:       listCertificatesCertificateBundleVersions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_id",
					Require: plugin.Optional,
				},
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreError: isNotFoundError([]string{"404"}),
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listCertificatesCertificateBundleVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	certificate := h.Item.(certificatesmanagement.CertificateSummary)
	certificateId := d.EqualsQualString("certificate_id")
	logger.Debug("oci_certificate_certificate_bundle_version.listCertificatesCertificateBundleVersions", "OCI_REGION", region)

	if certificateId != "" && certificateId != *certificate.Id {
		return nil, nil
	}

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_certificate_certificate_bundle_version.listCertificatesCertificateBundleVersions", "connection_error", err)
		return nil, err
	}
	request := certificates.ListCertificateBundleVersionsRequest{}
	request.CertificateId = certificate.Id
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	// The ListCertificateBundleVersions API don't support paggination
	response, err := session.CertificatesClient.ListCertificateBundleVersions(ctx, request)
	if err != nil {
		plugin.Logger(ctx).Error("oci_certificate_certificate_bundle_version.listCertificatesCertificateBundleVersions", "api_error", err)
		return nil, err
	}
	for _, respItem := range response.Items {
		d.StreamListItem(ctx, respItem)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getCertificatesCertificateBundleVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_certificate_certificate_bundle_version.getCertificatesCertificateBundle", "OCI_REGION", region)

	if h.Item == nil && !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	request := buildGetCertificatesCertificateBundleFilters(d.EqualsQuals)

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("oci_certificate_certificate_bundle_version.getCertificatesCertificateBundle", "connection_error", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesClient.GetCertificateBundle(ctx, request)
	if err != nil {
		logger.Error("oci_certificate_certificate_bundle_version.getCertificatesCertificateBundle", "api_error", err)
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

	return request
}
