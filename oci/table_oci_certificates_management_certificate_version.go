package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

// TABLE DEFINITION
func tableCertificatesManagementCertificateVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificates_management_certificate_version",
		Description:      "OCI Certificate Version",
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
			Hydrate: getCertificatesManagementCertificateVersion,
		},
		List: &plugin.ListConfig{
			Hydrate: listCertificatesManagementCertificateVersions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "certificate_id",
					Require: plugin.Required,
				},
				{
					Name:    "version_number",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "certificate_id",
				Description: "The OCID of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "version_number",
				Description: "The version number of the certificate.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "stages",
				Description: "A list of stages of this entity.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "serial_number",
				Description: "A unique certificate identifier used in certificate revocation tracking, formatted as octets.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "issuer_ca_version_number",
				Description: "The version number of the issuing certificate authority (CA).",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "version_name",
				Description: "The name of the certificate version. When the value is not null, a name is unique across versions of a given certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "subject_alternative_names",
				Description: "A list of subject alternative names.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the certificate version, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "validity",
				Description: "TBC",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "revocation_status",
				Description: "TBC",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "time_created",
				Description: "Time that the Certificate Version was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Description: ColumnDescriptionTenant,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

// LIST FUNCTION
func listCertificatesManagementCertificateVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("listCertificatesManagementCertificateVersions", "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals
	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	//Build request parameters
	request := buildListCertificatesManagementCertificateVersionFilters(equalQuals)
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
		response, err := session.CertificatesManagementClient.ListCertificateVersions(ctx, request)
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
func getCertificatesManagementCertificateVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getCertificatesManagementCertificateVersion", "OCI_REGION", region)

	request := buildGetCertificatesManagementCertificateVersionFilters(d.KeyColumnQuals, h)

	// Create Session
	session, err := certificatesManagementService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesManagementCertificateVersion", "error_CertificatesManagementService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesManagementClient.GetCertificateVersion(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.CertificateVersion, nil
}

// Build additional list filters
func buildListCertificatesManagementCertificateVersionFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificatesmanagement.ListCertificateVersionsRequest {
	request := certificatesmanagement.ListCertificateVersionsRequest{}

	if equalQuals["certificate_id"] != nil {
		request.CertificateId = types.String(equalQuals["certificate_id"].GetStringValue())
	}

	if equalQuals["version_number"] != nil {
		request.VersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}

	return request
}

// Build additional filters
func buildGetCertificatesManagementCertificateVersionFilters(equalQuals plugin.KeyColumnEqualsQualMap, h *plugin.HydrateData) certificatesmanagement.GetCertificateVersionRequest {
	request := certificatesmanagement.GetCertificateVersionRequest{}

	if h.Item != nil {
		request.CertificateId = h.Item.(certificatesmanagement.CertificateVersionSummary).CertificateId
		request.CertificateVersionNumber = h.Item.(certificatesmanagement.CertificateVersionSummary).VersionNumber
	} else {
		request.CertificateId = types.String(equalQuals["certificate_id"].GetStringValue())
		request.CertificateVersionNumber = types.Int64(equalQuals["version_number"].GetInt64Value())
	}

	return request
}
