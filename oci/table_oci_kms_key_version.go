package oci

import (
	"context"
	"sync"

	"github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/keymanagement"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableKmsKeyVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_kms_key_version",
		Description:      "OCI KMS Key Version",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			ParentHydrate: listKmsVaults,
			Hydrate:       listKmsKeyVersions,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
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
				Name:        "key_name",
				Description: "A user-friendly name of the key. Does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the vault that contains this key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_name",
				Description: "The display name of the vault that contains this key version.",
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// Key version info

type KeyVersionInfo struct {
	keymanagement.KeyVersionSummary
	ManagementEndpoint string
	KeyName            string
	VaultName          string
}

//// LIST FUNCTION

func listKmsKeyVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)

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

	logger.Debug("listKmsKeyVersions", "OCI_REGION", region, "Compartment", compartment, "Vault Name", *vaultData.DisplayName)

	// Create Session
	session, err := kmsManagementService(ctx, d, region, *vaultData.ManagementEndpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListKeysRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	var keyList []keymanagement.KeySummary

	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsManagementClient.ListKeys(ctx, request)
		if err != nil {
			return nil, err
		}

		keyList = append(keyList, response.Items...)

		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	var wg sync.WaitGroup
	keyVersionCh := make(chan []KeyVersionInfo, len(keyList))
	errorCh := make(chan error, len(keyList))

	// Iterating all the available keys
	for _, item := range keyList {
		wg.Add(1)
		go getRowDataForKeyVersionAsync(ctx, item, session, &wg, keyVersionCh, *vaultData.ManagementEndpoint, *vaultData.DisplayName, errorCh)
	}

	// wait for all keys to be processed
	wg.Wait()

	// NOTE: close channel before ranging over results
	close(keyVersionCh)
	close(errorCh)

	for err := range errorCh {
		// return the first error
		return nil, err
	}

	for keyVersion := range keyVersionCh {
		for _, version := range keyVersion {
			d.StreamListItem(ctx, version)
		}
	}

	return nil, err
}

func getRowDataForKeyVersionAsync(ctx context.Context, item keymanagement.KeySummary, session *session, wg *sync.WaitGroup, versionCh chan []KeyVersionInfo, endpoint string, vaultName string, errorCh chan error) {
	defer wg.Done()

	rowData, err := getRowDataForKeyVersion(ctx, item, session, endpoint, vaultName)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		versionCh <- rowData
	}
}

// List all the available key versions
func getRowDataForKeyVersion(ctx context.Context, item keymanagement.KeySummary, session *session, endpoint string, vaultName string) ([]KeyVersionInfo, error) {
	request := keymanagement.ListKeyVersionsRequest{
		KeyId: item.Id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	var versionInfo []KeyVersionInfo

	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsManagementClient.ListKeyVersions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, keyVersion := range response.Items {
			versionInfo = append(versionInfo, KeyVersionInfo{keyVersion, endpoint, *item.DisplayName, vaultName})
		}

		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return versionInfo, nil
}

//// HYDRATE FUNCTION

func getKmsKeyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKeyVersion")

	keyVersion := h.Item.(KeyVersionInfo)
	endpoint := keyVersion.ManagementEndpoint
	region := ociRegionNameFromId(*keyVersion.Id)

	// Create Session
	session, err := kmsManagementService(ctx, d, string(region), endpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.GetKeyVersionRequest{
		KeyId:        keyVersion.KeyId,
		KeyVersionId: keyVersion.Id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.KmsManagementClient.GetKeyVersion(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.KeyVersion, nil
}
