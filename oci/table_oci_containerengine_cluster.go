package oci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOciContainerEngineCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_containerengine_cluster",
		Description: "OCI Container Engine Cluster",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           getContainerEngineCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerEngineClusters,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
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
				Name:        "name",
				Description: "A user-friendly name. It does not have to be unique, and it is changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The state of the cluster masters.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "image_policy_config_enabled",
				Description: "Whether the image verification policy is enabled. Defaults to false. If set to true, the images will be verified against the policy at runtime.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ImagePolicyConfig.IsPolicyEnabled"),
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the KMS key to be used as the master encryption key for Kubernetes secret encryption.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getContainerEngineCluster,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "kubernetes_version",
				Description: "The version of Kubernetes running on the cluster masters.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current 'lifecycleState'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vcn_id",
				Description: "The OCID of the virtual cloud network (VCN) in which the cluster exists.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},

			//json fields
			{
				Name:        "available_kubernetes_upgrades",
				Description: "Available Kubernetes versions to which the clusters masters may be upgraded.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "endpoints",
				Description: "Endpoints served up by the cluster masters.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "endpoint_config",
				Description: "The network configuration for access to the Cluster control plane.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata",
				Description: "Metadata about the cluster.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "options",
				Description: "Optional attributes for the cluster.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
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
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listContainerEngineClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listContainerEngineClusters", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := containerEngineService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request, isValid := buildContainerEngineClusterFilters(equalQuals, logger)
	if !isValid {
		return nil, nil
	}
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ContainerEngineClient.ListClusters(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, cluster := range response.Items {
			d.StreamListItem(ctx, cluster)

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

func getContainerEngineCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getContainerEngineClusters", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(containerengine.ClusterSummary).Id
	} else {

		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
		id = d.EqualsQuals["id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := containerEngineService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := containerengine.GetClusterRequest{
		ClusterId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ContainerEngineClient.GetCluster(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Cluster, nil
}

// Build additional filters
func buildContainerEngineClusterFilters(equalQuals plugin.KeyColumnEqualsQualMap, logger hclog.Logger) (containerengine.ListClustersRequest, bool) {
	request := containerengine.ListClustersRequest{}
	isValid := true

	if equalQuals["name"] != nil && strings.Trim(equalQuals["name"].GetStringValue(), " ") != "" {
		request.Name = types.String(equalQuals["name"].GetStringValue())
	}
	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		if isValidContainerEngineClusterLifecycleState(lifecycleState) {
			request.LifecycleState = []containerengine.ClusterLifecycleStateEnum{containerengine.ClusterLifecycleStateEnum(lifecycleState)}
		} else {
			isValid = false
		}
	}
	return request, isValid
}

func isValidContainerEngineClusterLifecycleState(state string) bool {
	stateType := containerengine.ClusterLifecycleStateEnum(state)
	switch stateType {
	case containerengine.ClusterLifecycleStateActive, containerengine.ClusterLifecycleStateCreating, containerengine.ClusterLifecycleStateDeleted, containerengine.ClusterLifecycleStateDeleting, containerengine.ClusterLifecycleStateFailed, containerengine.ClusterLifecycleStateUpdating:
		return true
	}
	return false
}
