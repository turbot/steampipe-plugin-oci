package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableKmsKeyVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_kms_key_version",
		Description:      "OCI KMS Key Version",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"key_id", "management_endpoint", "region"}),
			Hydrate:    listKmsKeyVersions,
		},
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_id",
				Description: "The OCID of the master encryption key associated with this key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the vault that contains this key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "management_endpoint",
				Description: "The service endpoint to perform management operations against.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The key version's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "origin",
				Description: "The source of the key material. When this value is INTERNAL, Key Management created the key material. When this value is EXTERNAL, the key material was imported from an external source.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key",
				Description: "The public key in PEM format which will be populated only in case of RSA and ECDSA keys.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsKeyVersion,
			},
			{
				Name:        "restored_from_key_version_id",
				Description: "The OCID of the key version from which this key version was restored.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsKeyVersion,
			},
			{
				Name:        "time_created",
				Description: "The date and time this key version was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property to indicate when to delete the key version.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tenant_id",
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// Key version info

type KeyVersionInfo struct {
	keymanagement.KeyVersionSummary
	ManagementEndpoint string
	Region             string
}

//// LIST FUNCTION

func listKmsKeyVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	keyId := d.EqualsQuals["key_id"].GetStringValue()
	endpoint := d.EqualsQuals["management_endpoint"].GetStringValue()
	region := d.EqualsQuals["region"].GetStringValue()

	// handle empty keyId, endpoint and region in list call
	if keyId == "" || endpoint == "" || region == "" {
		return nil, nil
	}

	// Create Session
	session, err := kmsManagementService(ctx, d, region, endpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListKeyVersionsRequest{
		KeyId: types.String(keyId),
		Limit: types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsManagementClient.ListKeyVersions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, keyVersion := range response.Items {
			d.StreamListItem(ctx, KeyVersionInfo{keyVersion, endpoint, region})

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

func getKmsKeyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKeyVersion")

	keyVersion := h.Item.(KeyVersionInfo)
	endpoint := keyVersion.ManagementEndpoint
	region := keyVersion.Region

	// Create Session
	session, err := kmsManagementService(ctx, d, region, endpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.GetKeyVersionRequest{
		KeyId:        keyVersion.KeyId,
		KeyVersionId: keyVersion.Id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.KmsManagementClient.GetKeyVersion(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.KeyVersion, nil
}
