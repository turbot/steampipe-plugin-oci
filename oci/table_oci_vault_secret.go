package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/vault"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableVaultSecret(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_vault_secret",
		Description: "OCI Vault Secret",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			Hydrate:           getVaultSecret,
			ShouldIgnoreError: isNotFoundError([]string{"400", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listVaultSecrets,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "vault_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecretName"),
			},
			{
				Name:        "id",
				Description: "The OCID of the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "key_id",
				Description: "The OCID of the master encryption key that is used to encrypt the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current lifecycle state of the secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the Vault in which the secret exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "current_version_number",
				Description: "The version number of the secret that's currently in use.",
				Hydrate:     getVaultSecret,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "description",
				Description: "A brief description of the secret.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the secret's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "A property indicating when the secret was created.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_of_current_version_expiry",
				Description: "An optional property indicating when the current secret version will expire.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeOfCurrentVersionExpiry.Time"),
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the secret.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "metadata",
				Description: "Additional metadata that you can use to provide context about how to use the secret or during rotation or other administrative tasks.",
				Hydrate:     getVaultSecret,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "secret_rules",
				Description: "A list of rules that control how the secret is used and managed.",
				Hydrate:     getVaultSecret,
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(vaultSecretTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SecretName"),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listVaultSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Trace("listVaultSecrets", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := vaultService(ctx, d, region)
	if err != nil {
		logger.Error("listVaultSecrets", "error_vaultService", err)
		return nil, err
	}

	// Build request parameters
	request := buildVaultSecretFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
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
		response, err := session.VaultClient.ListSecrets(ctx, request)
		if err != nil {
			logger.Error("listVaultSecrets", "error_ListSecrets", err)
			if ociErr, ok := err.(common.ServiceError); ok {
				if ociErr.GetCode() == "InvalidParameter" || ociErr.GetCode() == "BadErrorResponse" {
					return nil, nil
				}
			}
			return nil, err
		}

		for _, vault := range response.Items {
			d.StreamListItem(ctx, vault)

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

	return nil, nil
}

//// HYDRATE FUNCTION

func getVaultSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getVaultSecret")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getVaultSecret", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		i := h.Item.(vault.SecretSummary)
		id = *i.Id
	} else {
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// handle empty secret id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := vaultService(ctx, d, region)
	if err != nil {
		logger.Error("getVaultSecret", "error_vaultService", err)
		return nil, err
	}

	request := vault.GetSecretRequest{
		SecretId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.VaultClient.GetSecret(ctx, request)
	if err != nil {
		logger.Error("getVaultSecret", "error_GetSecret", err)
		return nil, err
	}

	return response.Secret, nil
}

//// TRANSFORM FUNCTION

func vaultSecretTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := vaultSecretFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := vaultSecretDefinedTags(d.HydrateItem)

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

func vaultSecretFreeformTags(item interface{}) map[string]string {
	switch item := item.(type) {
	case vault.Secret:
		return item.FreeformTags
	case vault.SecretSummary:
		return item.FreeformTags
	}
	return nil
}

func vaultSecretDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item := item.(type) {
	case vault.Secret:
		return item.DefinedTags
	case vault.SecretSummary:
		return item.DefinedTags
	}
	return nil
}

// Build additional filters
func buildVaultSecretFilters(equalQuals plugin.KeyColumnEqualsQualMap) vault.ListSecretsRequest {
	request := vault.ListSecretsRequest{}

	if equalQuals["name"] != nil {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}
	if equalQuals["vault_id"] != nil {
		request.VaultId = types.String(equalQuals["vault_id"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = vault.SecretSummaryLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	return request
}
