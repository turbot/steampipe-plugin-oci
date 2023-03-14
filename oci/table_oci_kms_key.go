package oci

import (
	"context"
	"strings"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
					Name:    "algorithm",
					Require: plugin.Optional,
				},
				{
					Name:    "curve_id",
					Require: plugin.Optional,
				},
				{
					Name:    "length",
					Require: plugin.Optional,
				},
				{
					Name:    "protection_mode",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
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
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	equalQuals := d.EqualsQuals

	vaultData := h.Item.(keymanagement.VaultSummary)

	// skip the API call if vault is any of the below state
	if helpers.StringSliceContains([]string{"CREATING", "DELETING", "DELETED", "RESTORING"}, types.ToString(vaultData.LifecycleState)) {
		return nil, nil
	}

	// skip the API call if vault region doesn't match matrix region
	if ociRegionNameFromId(*vaultData.Id) != common.StringToRegion(region) {
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

	var wg sync.WaitGroup
	compartments, err := listAllCompartments(ctx, d)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildKmsKeyFilters(equalQuals)
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

	// List All keys available in various compartments other than vault compartment
	errorCh := make(chan error, len(compartments))
	for _, compartment := range compartments {
		request.CompartmentId = compartment.Id
		wg.Add(1)
		go getKmsKeyAsync(ctx, d, request, session, &wg, vaultData, errorCh)
	}

	// wait for all keys to be processed
	wg.Wait()

	// NOTE: close channel before ranging over error
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	return nil, nil
}

func getKmsKeyAsync(ctx context.Context, d *plugin.QueryData, request keymanagement.ListKeysRequest, session *session, wg *sync.WaitGroup, vaultData keymanagement.VaultSummary, errorCh chan error) {
	defer wg.Done()
	err := getKmsKeyAsyncData(ctx, d, request, session, vaultData)
	if err != nil {
		errorCh <- err
	}
}
func getKmsKeyAsyncData(ctx context.Context, d *plugin.QueryData, request keymanagement.ListKeysRequest, session *session, vaultData keymanagement.VaultSummary) error {
	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsManagementClient.ListKeys(ctx, request)
		if err != nil {
			return err
		}

		for _, key := range response.Items {
			d.StreamListItem(ctx, KeyInfo{key, *vaultData.ManagementEndpoint, *vaultData.DisplayName})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil
			}
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil
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
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
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
func ociRegionNameFromId(resourceId string) common.Region {
	id := types.SafeString(resourceId)
	splittedID := strings.Split(id, ".")
	regionName := common.StringToRegion(types.SafeString(splittedID[3]))
	return regionName
}

// Build additional filters
func buildKmsKeyFilters(equalQuals plugin.KeyColumnEqualsQualMap) keymanagement.ListKeysRequest {
	request := keymanagement.ListKeysRequest{}

	if equalQuals["algorithm"] != nil {
		request.Algorithm = keymanagement.ListKeysAlgorithmEnum(equalQuals["algorithm"].GetStringValue())
	}
	if equalQuals["curve_id"] != nil {
		request.CurveId = keymanagement.ListKeysCurveIdEnum(equalQuals["curve_id"].GetStringValue())
	}
	if equalQuals["length"] != nil {
		request.Length = types.Int(int(equalQuals["length"].GetInt64Value()))
	}
	if equalQuals["protection_mode"] != nil {
		request.ProtectionMode = keymanagement.ListKeysProtectionModeEnum(equalQuals["protection_mode"].GetStringValue())
	}

	return request
}
