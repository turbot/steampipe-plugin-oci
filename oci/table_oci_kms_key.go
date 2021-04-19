package oci

import (
	"context"

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
		List: &plugin.ListConfig{
			ParentHydrate: listKmsVaults,
			Hydrate:       listKmsKeys,
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
		},
	}
}

//// LIST FUNCTION

func listKmsKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listKmsKeys", "Compartment", compartment, "OCI_REGION", region)

	vaultData := h.Item.(*keymanagement.VaultSummary)
	endpoint := *vaultData.CryptoEndpoint

	// Create Session
	session, err := kmsManagementService(ctx, d, region, endpoint)

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
			d.StreamListItem(ctx, key)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}
