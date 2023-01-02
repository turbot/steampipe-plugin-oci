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
func tableCertificatesCaBundle(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_certificates_ca_bundle",
		Description:      "OCI Ca Bundle",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("ca_bundle_id"),
			Hydrate:    getCertificatesCaBundle,
		},
		List: &plugin.ListConfig{
			Hydrate:    getCertificatesCaBundle,
			KeyColumns: plugin.SingleColumn("ca_bundle_id"),
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "ca_bundle_id",
				Description: "The OCID of the CA bundle.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "name",
				Description: "A user-friendly name for the CA bundle. Names are unique within a compartment. Valid characters include uppercase or lowercase letters, numbers, hyphens, underscores, and periods.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ca_bundle_pem",
				Description: "Certificates (in PEM format) in the CA bundle. Can be of arbitrary length.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
func getCertificatesCaBundle(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getCertificatesCaBundle", "OCI_REGION", region, "keyQuals", d.KeyColumnQuals)

	request := buildGetCertificatesCaBundleFilters(d.KeyColumnQuals)

	// Create Session
	session, err := certificatesService(ctx, d, region)
	if err != nil {
		logger.Error("getCertificatesCaBundle", "error_CertificatesService", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.CertificatesClient.GetCaBundle(ctx, request)

	if err != nil {
		return nil, err
	}
	return response.CaBundle, nil
}

// Build additional filters
func buildGetCertificatesCaBundleFilters(equalQuals plugin.KeyColumnEqualsQualMap) certificates.GetCaBundleRequest {
	request := certificates.GetCaBundleRequest{}

	if equalQuals["ca_bundle_id"] != nil && equalQuals["ca_bundle_id"].GetStringValue() != "" {
		request.CaBundleId = types.String(equalQuals["ca_bundle_id"].GetStringValue())
	}

	return request
}
