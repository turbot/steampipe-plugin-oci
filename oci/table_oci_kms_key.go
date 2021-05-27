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

//// TABLE DEFINITION

func tableKmsKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_kms_key",
		Description: "OCI KMS Key",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.AllColumns([]string{"management_endpoint", "id"}),
		// 	Hydrate:    getKmsKey,
		// },
		List: &plugin.ListConfig{
			Hydrate: listKmsKeys,
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
				Description: "The OCID of the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vault_id",
				Description: "The OCID of the vault that contains the key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
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
				Name:        "time_created",
				Description: "The date and time the key was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "algorithm",
				Description: "The algorithm used by a key's key versions to encrypt or decrypt.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protection_mode",
				Description: "The key's protection mode indicates how the key persists and where cryptographic operations that use the key are performed.",
				Type:        proto.ColumnType_STRING,
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
			// {
			// 	Name:        "tags",
			// 	Description: ColumnDescriptionTags,
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.From(vaultTags),
			// },
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
				Hydrate:     getTenantId,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listKmsKeys", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := kmsVaultService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := keymanagement.ListVaultsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.KmsVaultClient.ListVaults(ctx, request)
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		errorCh := make(chan error, len(response.Items))
		for _, vault := range response.Items {
			if vault.LifecycleState != "ACTIVE" {
				return nil, nil
			}
			wg.Add(1)
			go getKmsKeyAsync(ctx, d, &wg, vault, errorCh)
		}
		// wait for all vaults to be processed
		wg.Wait()

		// NOTE: close channel before ranging over results
		close(errorCh)

		for err := range errorCh {
			// return the first error
			return nil, err
		}

		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

func getKmsKeyAsync(ctx context.Context, d *plugin.QueryData, wg *sync.WaitGroup, vault keymanagement.VaultSummary, errorCh chan error) {
	defer wg.Done()

	err := getKmsKeyDetails(ctx, d, vault)
	if err != nil {
		errorCh <- err
	}
}

func getKmsKeyDetails(ctx context.Context, d *plugin.QueryData, vault keymanagement.VaultSummary) error {
	compartment := *vault.CompartmentId
	endpoint := *vault.ManagementEndpoint
	region := strings.Split(endpoint, ".")[2]

	// Create Session
	session, err := kmsManagementService(ctx, d, region, endpoint)

	if err != nil {
		return err
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
			return err
		}

		for _, key := range response.Items {
			d.StreamListItem(ctx, key)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return err
}

//// HYDRATE FUNCTION

// func getKmsKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("getKmsKey")
// 	logger := plugin.Logger(ctx)
// 	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
// 	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
// 	logger.Debug("oci.getKmsKey", "Compartment", compartment, "OCI_REGION", region)

// 	// Restrict the api call to only root compartment/ per region
// 	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
// 		return nil, nil
// 	}

// 	endpoint := d.KeyColumnQuals["management_endpoint"].GetStringValue()
// 	id := d.KeyColumnQuals["id"].GetStringValue()

// 	// handle empty key id in get call
// 	if strings.TrimSpace(id) == "" {
// 		return nil, nil
// 	}

// 	// Create Session
// 	session, err := kmsManagementService(ctx, d, region, endpoint)
// 	if err != nil {
// 		return nil, err
// 	}

// 	request := keymanagement.GetKeyRequest{
// 		KeyId: types.String(id),
// 		RequestMetadata: oci_common.RequestMetadata{
// 			RetryPolicy: getDefaultRetryPolicy(),
// 		},
// 	}

// 	response, err := session.KmsManagementClient.GetKey(ctx, request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response.Key, nil
// }
