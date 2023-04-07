package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreClusterNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_cluster_network",
		Description: "OCI Core Cluster Network",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"400"}),
			Hydrate:           getClusterNetwork,
		},
		List: &plugin.ListConfig{
			Hydrate: listClusterNetworks,
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
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the cluster network.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the cluster network.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the resource was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the resource was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},

			// json fields
			{
				Name:        "instance_pools",
				Description: "The instance pools in the cluster network.",
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
				Transform:   transform.From(clusterNetworkTags),
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
		}),
	}
}

//// LIST FUNCTION

func listClusterNetworks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeManagementService(ctx, d, region)
	if err != nil {
		logger.Error("oci_core_cluster_network.ListClusterNetworks", "connection_error", err)
		return nil, err
	}

	request := core.ListClusterNetworksRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}
	displayName := d.EqualsQualString("display_name")
	if displayName != "" {
		request.DisplayName = &displayName
	}
	lifecycleState := d.EqualsQualString("lifecycle_state")
	if lifecycleState != "" {
		request.DisplayName = &lifecycleState
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeManagementClient.ListClusterNetworks(ctx, request)
		if err != nil {
			logger.Error("oci_core_cluster_network.ListClusterNetworks", "api_error", err)
			return nil, err
		}

		for _, computeManagement := range response.Items {
			d.StreamListItem(ctx, computeManagement)

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

func getClusterNetwork(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	matrixRegion := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

	// For the us-phoenix-1 and us-ashburn-1 regions, `phx` and `iad` are returned by ListInstances api, respectively.
	// For all other regions, the full region name is returned.
	region := common.StringToRegion(types.SafeString(strings.Split(id, ".")[3]))

	// handle empty id and region check in get call
	if id == "" || region != common.StringToRegion(matrixRegion) {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeManagementService(ctx, d, matrixRegion)
	if err != nil {
		logger.Error("oci_core_cluster_network.getClusterNetwork", "connection_error", err)
		return nil, err
	}

	request := core.GetClusterNetworkRequest{
		ClusterNetworkId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ComputeManagementClient.GetClusterNetwork(ctx, request)
	if err != nil {
		logger.Error("oci_core_cluster_network.getClusterNetwork", "api_error", err)
		return nil, err
	}

	return response, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func clusterNetworkTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var tags map[string]interface{}
	switch item := d.HydrateItem.(type) {
	case core.ClusterNetworkSummary:
		if item.FreeformTags != nil {
			tags = map[string]interface{}{}
			for k, v := range item.FreeformTags {
				tags[k] = v
			}
		}

		if item.DefinedTags != nil {
			if tags == nil {
				tags = map[string]interface{}{}
			}
			for _, v := range item.DefinedTags {
				for key, value := range v {
					tags[key] = value
				}

			}
		}
	case core.GetClusterNetworkResponse:
		if item.FreeformTags != nil {
			tags = map[string]interface{}{}
			for k, v := range item.FreeformTags {
				tags[k] = v
			}
		}

		if item.DefinedTags != nil {
			if tags == nil {
				tags = map[string]interface{}{}
			}
			for _, v := range item.DefinedTags {
				for key, value := range v {
					tags[key] = value
				}

			}
		}
	}

	return tags, nil
}
