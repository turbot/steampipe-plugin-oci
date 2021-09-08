package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/keymanagement"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableKmsKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_kms_key",
		Description:      "OCI KMS Key",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			ParentHydrate: listKmsVaults,
			Hydrate:       listKmsKeys,
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
				Name:        "name",
				Description: "A user-friendly name of the key. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "id",
				Description: "The OCID of the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the vault that contains the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_name",
				Description: "The display name of the vault that contains the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "management_endpoint",
				Description: "The service endpoint to perform management operations against.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The key's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "algorithm",
				Description: "The algorithm used by a key's key versions to encrypt or decrypt.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "current_key_version",
				Description: "The OCID of the key version used in cryptographic operations.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsKey,
			},
			{
				Name:        "curve_id",
				Description: "Supported curve Ids for ECDSA keys.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsKey,
				Transform:   transform.FromField("KeyShape.CurveId"),
			},
			{
				Name:        "length",
				Description: "The length of the key.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getKmsKey,
				Transform:   transform.FromField("KeyShape.Length"),
			},
			{
				Name:        "protection_mode",
				Description: "The key's protection mode indicates how the key persists and where cryptographic operations that use the key are performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "restored_from_key_id",
				Description: "The OCID of the key from which this key was restored.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getKmsKey,
			},
			{
				Name:        "time_created",
				Description: "The date and time the key was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the key.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
				Hydrate:     getKmsKey,
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
				Transform:   transform.From(keyTags),
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

//// Key info

type KeyInfo struct {
	keymanagement.KeySummary
	ManagementEndpoint string
	VaultName          string
}

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	vaultData := h.Item.(keymanagement.VaultSummary)

	// skip the API call if vault is any of the below state
	if helpers.StringSliceContains([]string{"CREATING", "DELETING", "DELETED", "RESTORING"}, types.ToString(vaultData.LifecycleState)) {
		return nil, nil
	}

	// skip the API call if vault region doesn't match matrix region
	if ociRegionNameFromId(*vaultData.Id) != oci_common.StringToRegion(region) {
		return nil, nil
	}

	// skip the API call if vault compartment doesn't match matrix compartment
	if *vaultData.CompartmentId != compartment {
		return nil, nil
	}

	logger.Debug("listKmsKeys", "OCI_REGION", region, "Compartment", compartment, "Vault Name", *vaultData.DisplayName)

	// Create Session
	session, err := kmsManagementService(ctx, d, region, *vaultData.ManagementEndpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListKeysRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsManagementClient.ListKeys(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, key := range response.Items {
			d.StreamListItem(ctx, KeyInfo{key, *vaultData.ManagementEndpoint, *vaultData.DisplayName})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				response.OpcNextPage = nil
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

func getKmsKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKey")

	key := h.Item.(KeyInfo)
	endpoint := key.ManagementEndpoint
	region := ociRegionNameFromId(*key.Id)

	// Create Session
	session, err := kmsManagementService(ctx, d, string(region), endpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.GetKeyRequest{
		KeyId: key.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.KmsManagementClient.GetKey(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Key, nil
}

//// TRANSFORM FUNCTION

func keyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := d.HydrateItem.(KeyInfo).FreeformTags

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := d.HydrateItem.(KeyInfo).DefinedTags

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

// Extract OCI region name from the resource id
func ociRegionNameFromId(resourceId string) oci_common.Region {
	id := types.SafeString(resourceId)
	splittedID := strings.Split(id, ".")
	regionName := oci_common.StringToRegion(types.SafeString(splittedID[3]))
	return regionName
}
