package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/networkloadbalancer"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreNetworkLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_network_load_balancer",
		Description: "OCI Core Network Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreNetworkLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreNetworkLoadBalancers,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the network load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the network load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_details",
				Description: "A message describing the current state in more detail.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the network load balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "time_updated",
				Description: "The date and time the network load balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "subnet_id",
				Description: "The subnet in which the network load balancer is spawned OCIDs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubnetId"),
			},
			{
				Name:        "is_private",
				Description: "Whether the network load balancer has a virtual cloud network-local (private) IP address.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_preserve_source_destination",
				Description: "When enabled, the skipSourceDestinationCheck parameter is automatically enabled on the load balancer VNIC.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "network_load_balancer_health",
				Description: "The overall health status of the network load balancer.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getCoreNetworkLoadBalancerHealth,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network_security_group_ids",
				Description: "An array of network security groups OCIDs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "listeners",
				Description: "Listeners associated with the network load balancer.",
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
			{
				Name:        "system_tags",
				Description: ColumnDescriptionSystemTags,
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(networkLoadBalancerTags),
			},
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

func listCoreNetworkLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreNetworkLoadBalancers", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := networkLoadBalancerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := networkloadbalancer.ListNetworkLoadBalancersRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.NetworkLoadBalancerClient.ListNetworkLoadBalancers(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, networkLoadBalancer := range response.Items {
			d.StreamListItem(ctx, networkLoadBalancer)
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

func getCoreNetworkLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreNetworkLoadBalancer")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreNetworkLoadBalancer", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty nlb id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := networkLoadBalancerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := networkloadbalancer.GetNetworkLoadBalancerRequest{
		NetworkLoadBalancerId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.NetworkLoadBalancerClient.GetNetworkLoadBalancer(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.NetworkLoadBalancer, nil
}

func getCoreNetworkLoadBalancerHealth(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreNetworkLoadBalancerHealth")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreNetworkLoadBalancerHealth", "Compartment", compartment, "OCI_REGION", region)

	var id string
	switch h.Item.(type) {
	case networkloadbalancer.NetworkLoadBalancerSummary:
		id = *h.Item.(networkloadbalancer.NetworkLoadBalancerSummary).Id
	case networkloadbalancer.NetworkLoadBalancer:
		id = *h.Item.(networkloadbalancer.NetworkLoadBalancer).Id
	}

	// Create Session
	session, err := networkLoadBalancerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := networkloadbalancer.GetNetworkLoadBalancerHealthRequest{
		NetworkLoadBalancerId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.NetworkLoadBalancerClient.GetNetworkLoadBalancerHealth(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.NetworkLoadBalancerHealth, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. System Tags
// 2. Defined Tags
// 3. Free-form tags
func networkLoadBalancerTags(_ context.Context, d *transform.TransformData) (interface{}, error) {

	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	var systemTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case networkloadbalancer.NetworkLoadBalancerSummary:
		networkLoadBalancer := d.HydrateItem.(networkloadbalancer.NetworkLoadBalancerSummary)
		freeformTags = networkLoadBalancer.FreeformTags
		definedTags = networkLoadBalancer.DefinedTags
		systemTags = networkLoadBalancer.SystemTags
	case networkloadbalancer.NetworkLoadBalancer:
		networkLoadBalancer := d.HydrateItem.(networkloadbalancer.NetworkLoadBalancer)
		freeformTags = networkLoadBalancer.FreeformTags
		definedTags = networkLoadBalancer.DefinedTags
		systemTags = networkLoadBalancer.SystemTags
	}

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

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

	if systemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range systemTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
