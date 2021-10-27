package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/keymanagement"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableKmsVault(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_kms_vault",
		Description: "OCI KMS Vault",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getKmsVault,
		},
		List: &plugin.ListConfig{
			Hydrate: listKmsVaults,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of a vault.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "A vault's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time a vault was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "crypto_endpoint",
				Description: "The service endpoint to perform cryptographic operations against.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "management_endpoint",
				Description: "The service endpoint to perform management operations against.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "restored_from_vault_id",
				Description: "The OCID of the vault from which this vault was restored, if it was restored from a backup file.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
				Hydrate:     getKmsVault,
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property to indicate when to delete the vault.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
				Hydrate:     getKmsVault,
			},
			{
				Name:        "vault_type",
				Description: "The type of vault. Each type of vault stores keys with different degrees of isolation and has different options and pricing.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "wrappingkey_id",
				Description: "The OCID of the vault's wrapping key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
				Hydrate:     getKmsVault,
			},

			// tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(vaultTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id").Transform(ociRegionName),
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CompartmentId"),
			},
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

//// LIST FUNCTION

func listKmsVaults(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listKmsVaults", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := kmsVaultService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListVaultsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
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
		response, err := session.KmsVaultClient.ListVaults(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, vault := range response.Items {
			d.StreamListItem(ctx, vault)

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

//// HYDRATE FUNCTION

func getKmsVault(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsVault")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("oci.getKmsVault", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	var id string
	if h.Item != nil {
		i := h.Item.(keymanagement.VaultSummary)
		id = *i.Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	// handle empty vault id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := kmsVaultService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := keymanagement.GetVaultRequest{
		VaultId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.KmsVaultClient.GetVault(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Vault, nil
}

//// TRANSFORM FUNCTION

func vaultTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := vaultFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := vaultDefinedTags(d.HydrateItem)

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

func vaultFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case keymanagement.Vault:
		return item.FreeformTags
	case keymanagement.VaultSummary:
		return item.FreeformTags
	}
	return nil
}

func vaultDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case keymanagement.Vault:
		return item.DefinedTags
	case keymanagement.VaultSummary:
		return item.DefinedTags
	}
	return nil
}
