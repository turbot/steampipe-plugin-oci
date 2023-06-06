package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesManagementCertificateAuthorityVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificate_management_certificate_authority_version",
		Description:      "OCI Certificate Authority Version",
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
			},
			Hydrate: getCertificatesManagementCertificateAuthorityVersion,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCertificatesManagementCertificateAuthorities,
			Hydrate:       listCertificatesManagementCertificateAuthorityVersions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_authority_id",
					Require: plugin.Optional,
				},
				{
					Name:    "version_number",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "certificate_authority_id",
				Description: "The OCID of the CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_number",
				Description: "The version number of this CA.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Authority Version was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the CA version.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "serial_number",
				Description: "A unique certificate identifier used in certificate revocation tracking, formatted as octets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "issuer_ca_version_number",
				Description: "The version number of the issuing CA.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "version_name",
				Description: "The name of the CA version. When the value is not null, a name is unique across versions for a given CA.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "stages",
				Description: "A list of rotation states for this CA version.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subject_alternative_names",
				Description: "A list of subject alternative names. A subject alternative name specifies the domain names, including subdomains, and IP addresses covered by the certificates issued by this CA.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCertificatesManagementCertificateAuthorityVersion,
			},
			{
				Name:        "validity",
				Description: "Certificate Authority validity details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "revocation_status",
				Description: "Revocation details for the CA.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionName"),
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

func listCertificatesManagementCertificateAuthorityVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	ca := h.Item.(certificatesmanagement.CertificateAuthoritySummary)
	logger.Debug("oci_certificate_management_certificate_authority_version.listCertificatesManagementCertificateAuthorityVersions", "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("oci_certificate_management_certificate_authority_version.listCertificatesManagementCertificateAuthorityVersions", "connection_error", err)
		return nil, err
	}

	if d.EqualsQualString("certificate_authority_id") != "" && d.EqualsQualString("certificate_authority_id") != *ca.Id {
		return nil, nil
	}

	//Build request parameters
	request := buildListCertificatesManagementCertificateAuthorityVersionFilters(equalQuals)
	request.CertificateAuthorityId = ca.Id
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
		response, err := session.CertificatesManagementClient.ListCertificateAuthorityVersions(ctx, request)
		if err != nil {
			plugin.Logger(ctx).Error("oci_certificate_management_certificate_authority_version.listCertificatesManagementCertificateAuthorityVersions", "api_error", err)
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

func getCertificatesManagementCertificateAuthorityVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("oci_certificate_management_certificate_authority_version.getCertificatesManagementCertificateAuthorityVersion", "OCI_REGION", region)

	request := buildGetCertificatesManagementCertificateAuthorityVersionFilters(d.EqualsQuals, h)

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		logger.Error("oci_certificate_management_certificate_authority_version.getCertificatesManagementCertificateAuthorityVersion", "connection_error", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesManagementClient.GetCertificateAuthorityVersion(ctx, request)
	if err != nil {
		logger.Error("oci_certificate_management_certificate_authority_version.getCertificatesManagementCertificateAuthorityVersion", "api_error", err)
		return nil, err
	}
	return response.CertificateAuthorityVersion, nil
}

// Build additional list filters
func buildListCertificatesManagementCertificateAuthorityVersionFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificatesmanagement.ListCertificateAuthorityVersionsRequest {
	request := certificatesmanagement.ListCertificateAuthorityVersionsRequest{}

	if equalQuals["version_number"] != nil {
		request.VersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}

	return request
}

// Build additional filters
func buildGetCertificatesManagementCertificateAuthorityVersionFilters(equalQuals plugin.KeyColumnEqualsQualMap, h *plugin.HydrateData) certificatesmanagement.GetCertificateAuthorityVersionRequest {
	request := certificatesmanagement.GetCertificateAuthorityVersionRequest{}

	if h.Item != nil {
		request.CertificateAuthorityId = h.Item.(certificatesmanagement.CertificateAuthorityVersionSummary).CertificateAuthorityId
		request.CertificateAuthorityVersionNumber = h.Item.(certificatesmanagement.CertificateAuthorityVersionSummary).VersionNumber
	} else {
		request.CertificateAuthorityId = types.String(equalQuals["certificate_authority_id"].GetStringValue())
		request.CertificateAuthorityVersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}

	return request
}
