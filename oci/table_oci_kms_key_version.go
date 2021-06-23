package oci

import (
	"context"
	"strings"
	"sync"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/keymanagement"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type keyVersionInfo struct {
	KeyVersion         keymanagement.KeyVersionSummary
	ManagementEndpoint string
}

//// TABLE DEFINITION

func tableKmsKeyVersion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_kms_key_version",
		Description:      "OCI KMS Key Version",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getKmsKeyVersion,
		},
		List: &plugin.ListConfig{
			// ParentHydrate: listKmsKeys,
			Hydrate: listKmsKeyVersions,
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
				Description: "The OCID of the key version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the vault that contains the key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The key's current lifecycle state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the key version was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_of_deletion",
				Description: "An optional property indicating when to delete the key version.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeOfDeletion.Time"),
			},
			{
				Name:        "origin",
				Description: "The source of the key material.",
				Type:        proto.ColumnType_STRING,
			},

			// Steampipe tandard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// OCI standard columns
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

//// LIST FUNCTION

func listKmsKeyVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account

	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	// compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	vaultData := h.Item.(keymanagement.VaultSummary)

	session, err := kmsManagementService(ctx, d, region, *vaultData.ManagementEndpoint)
	if err != nil {
		return nil, err
	}

	var kmsKeys []KeyInfo

	// List kms keys
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
			kmsKeys = append(kmsKeys, KeyInfo{key, *vaultData.ManagementEndpoint, *vaultData.DisplayName})
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	var wg sync.WaitGroup
	keyVersionCh := make(chan []keyVersionInfo, len(kmsKeys))
	errorCh := make(chan error, len(kmsKeys))

	// Iterating all the available keys
	for _, item := range kmsKeys {
		wg.Add(1)
		go getRowDataForKeyVersionAsync(ctx, d, item, &wg, keyVersionCh, errorCh, region)
	}

	// wait for all keys to be processed
	wg.Wait()
	close(keyVersionCh)
	close(errorCh)

	for err := range errorCh {
		return nil, err
	}

	for item := range keyVersionCh {
		for _, data := range item {
			d.StreamLeafListItem(ctx, keyVersionInfo{data.KeyVersion, data.ManagementEndpoint})
		}
	}

	return nil, err
}

func getRowDataForKeyVersionAsync(ctx context.Context, d *plugin.QueryData, item KeyInfo, wg *sync.WaitGroup, subnetCh chan []keyVersionInfo, errorCh chan error, region string) {
	defer wg.Done()

	rowData, err := getRowDataForKeyVersion(ctx, d, item, region)
	if err != nil {
		errorCh <- err
	} else if rowData != nil {
		subnetCh <- rowData
	}
}

func getRowDataForKeyVersion(ctx context.Context, d *plugin.QueryData, item KeyInfo, region string) ([]keyVersionInfo, error) {

	session, err := kmsManagementService(ctx, d, region, item.ManagementEndpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListKeyVersionsRequest{
		KeyId: item.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	var items []keyVersionInfo

	listKeyVersion, err := session.KmsManagementClient.ListKeyVersions(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, keyVersion := range listKeyVersion.Items {
		items = append(items, keyVersionInfo{keyVersion, item.ManagementEndpoint})
	}

	return items, err
}

//// HYDRATE FUNCTION

func getKmsKeyVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKmsKeyVersion")

	key := h.Item.(KeyInfo)
	endpoint := key.ManagementEndpoint
	region := ociRegionNameFromKeyVersionId(*key.Id)

	// Create Session
	session, err := kmsManagementService(ctx, d, string(region), endpoint)
	if err != nil {
		return nil, err
	}

	request := keymanagement.GetKeyVersionRequest{
		KeyId: key.Id,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.KmsManagementClient.GetKeyVersion(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.KeyVersion, nil
}

//// TRANSFORM FUNCTION

// Extract OCI region name from the resource id
func ociRegionNameFromKeyVersionId(resourceId string) oci_common.Region {
	id := types.SafeString(resourceId)
	splittedID := strings.Split(id, ".")
	regionName := oci_common.StringToRegion(types.SafeString(splittedID[3]))
	return regionName
}
