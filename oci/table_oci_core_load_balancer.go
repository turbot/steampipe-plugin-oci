package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreLoadBalancer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_load_balancer",
		Description: "OCI Core Load Balancer",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreLoadBalancer,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreLoadBalancers,
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
				Description: "The OCID of the load balancer.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name of the load balancer.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The load balancer's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the load balancer was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "shape_name",
				Description: "A template that determines the total pre-provisioned bandwidth (ingress plus egress).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "is_private",
				Description: "Whether the load balancer has a VCN-local (private) IP address.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "backend_sets",
				Description: "The configuration of a load balancer backend set.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "certificates",
				Description: "The configuration details of a certificate bundle.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "hostnames",
				Description: "A hostname resource associated with a load balancer for use by one or more listeners.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ip_addresses",
				Description: "An array of IP addresses.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "listeners",
				Description: "The listener's configuration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "network_security_group_ids",
				Description: "An array of NSG OCIDs associated with the load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "path_route_sets",
				Description: "A named set of path route rules.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "routing_policies",
				Description: "A named ordered list of routing rules that is applied to a listener.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "rule_sets",
				Description: "A named set of rules associated with a load balancer.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "shape_details",
				Description: "The configuration details to update load balancer to a different shape.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ssl_cipher_suites",
				Description: "The configuration details of an SSL cipher suite.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "subnet_ids",
				Description: "An array of subnet OCIDs.",
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
				Description: "System tags for this resource.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(loadBalancerTags),
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

func listCoreLoadBalancers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.listCoreLoadBalancers", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := loadBalancerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := loadbalancer.ListLoadBalancersRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int64(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["display_name"] != nil {
		displayName := equalQuals["display_name"].GetStringValue()
		request.DisplayName = types.String(displayName)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.LifecycleState = loadbalancer.LoadBalancerLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int64(int64(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.LoadBalancerClient.ListLoadBalancers(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, loadBalancer := range response.Items {
			d.StreamListItem(ctx, loadBalancer)

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

func getCoreLoadBalancer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreLoadBalancer")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci.getCoreLoadBalancer", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

	// handle empty id in get call
	if strings.TrimSpace(id) == "" {
		return nil, nil
	}

	// Create Session
	session, err := loadBalancerService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := loadbalancer.GetLoadBalancerRequest{
		LoadBalancerId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.LoadBalancerClient.GetLoadBalancer(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.LoadBalancer, nil
}

//// TRANSFORM FUNCTION

func loadBalancerTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	loadBalancer := d.HydrateItem.(loadbalancer.LoadBalancer)

	var tags map[string]interface{}

	if loadBalancer.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range loadBalancer.FreeformTags {
			tags[k] = v
		}
	}

	if loadBalancer.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range loadBalancer.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}
		}
	}

	if loadBalancer.SystemTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range loadBalancer.SystemTags {
			for key, value := range v {
				tags[key] = value
			}
		}
	}

	return tags, nil
}
