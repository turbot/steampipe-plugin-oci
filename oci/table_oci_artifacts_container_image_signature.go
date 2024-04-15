package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/artifacts"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableArtifactContainerImageSignature(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_artifacts_container_image_signature",
		Description:      "OCI Artifacts Container Image Signature",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getArtifactContainerImageSignature,
		},
		List: &plugin.ListConfig{
			Hydrate: listArtifactContainerImageSignatures,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "kms_key_id",
					Require: plugin.Optional,
				},
				{
					Name:    "kms_key_version_id",
					Require: plugin.Optional,
				},
				{
					Name:    "signing_algorithm",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "The last 10 characters of the kmsKeyId, the last 10 characters of the kmsKeyVersionId, the signingAlgorithm, and the last 10 characters of the signatureId.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the container image signature.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that Container Image Signature was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "created_by",
				Description: "The id of the user or principal that created the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getArtifactContainerImageSignature,
			},
			{
				Name:        "image_id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the container image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the kmsKeyId used to sign the container image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "kms_key_version_id",
				Description: "The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the kmsKeyVersionId used to sign the container image.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "message",
				Description: "The base64 encoded signature payload that was signed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signature",
				Description: "The signature of the message field using the kmsKeyId, the kmsKeyVersionId, and the signingAlgorithm.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "signing_algorithm",
				Description: "The algorithm to be used for signing. These are the only supported signing algorithms for container images.",
				Type:        proto.ColumnType_STRING,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,

				Transform: transform.FromField("DisplayName"),
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

func listArtifactContainerImageSignatures(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifacts_container_image_signature.listArtifactContainerImageSignatures", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_container_image_signature.listArtifactContainerImageSignatures", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildArtifactContainerImageSignatureFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(100)
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
		response, err := session.ArtifactClient.ListContainerImageSignatures(ctx, request)
		if err != nil {
			logger.Error("oci_artifacts_container_image_signature.listArtifactContainerImageSignatures", "api_error", err)
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

func getArtifactContainerImageSignature(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_artifacts_container_image_signature.getArtifactContainerImageSignature", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(artifacts.ContainerImageSignatureSummary).Id
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session

	session, err := artifactService(ctx, d, region)
	if err != nil {
		logger.Error("oci_artifacts_container_image_signature.getArtifactContainerImageSignature", "connection_error", err)
		return nil, err
	}

	request := artifacts.GetContainerImageSignatureRequest{
		ImageSignatureId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ArtifactClient.GetContainerImageSignature(ctx, request)
	if err != nil {
		logger.Error("oci_artifacts_container_image_signature.getArtifactContainerImageSignature", "api_error", err)
		return nil, err
	}
	return response.ContainerImageSignature, nil
}

// Build additional filters
func buildArtifactContainerImageSignatureFilters(equalQuals plugin.KeyColumnEqualsQualMap) artifacts.ListContainerImageSignaturesRequest {
	request := artifacts.ListContainerImageSignaturesRequest{}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	if equalQuals["kms_key_id"] != nil {
		request.KmsKeyId = types.String(equalQuals["kms_key_id"].GetStringValue())
	}

	if equalQuals["kms_key_version_id"] != nil {
		request.KmsKeyVersionId = types.String(equalQuals["kms_key_version_id"].GetStringValue())
	}

	if equalQuals["signing_algorithm"] != nil {
		request.SigningAlgorithm = artifacts.ListContainerImageSignaturesSigningAlgorithmEnum(equalQuals["signing_algorithm"].GetStringValue())
	}

	return request
}
