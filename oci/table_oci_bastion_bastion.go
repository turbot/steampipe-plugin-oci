package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	"strings"
)

//// TABLE DEFINITION

func tableBastion(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_bastion_bastion",
		Description:      "OCI Bastion Bastion",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBastion,
		},
		List: &plugin.ListConfig{
			Hydrate: listBastions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the bastion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The display name of the bastion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bastion_type",
				Description: "The type of bastion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dns_proxy_status",
				Description: "The current DNS proxy status of the bastion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "client_cidr_block_allow_list",
				Description: "A list of address ranges in CIDR notation that you want to allow to connect to sessions hosted by this bastion.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "max_session_ttl_in_seconds",
				Description: "The maximum amount of time that any session on the bastion can remain active.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_sessions_allowed",
				Description: "The maximum number of active sessions allowed on the bastion.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "private_endpoint_ip_address",
				Description: "The private IP address of the created private endpoint.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_IPADDR,
			},
			{
				Name:        "static_jump_host_ip_addresses",
				Description: "A list of IP addresses of the hosts that the bastion has access to. Not applicable to standard bastions.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "phone_book_entry",
				Description: "The phonebook entry of the customer's team, which can't be changed after creation. Not applicable to standard bastions.",
				Hydrate:     getBastion,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_vcn_id",
				Description: "The unique identifier (OCID) of the virtual cloud network (VCN) that the bastion connects to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_subnet_id",
				Description: "The unique identifier (OCID) of the subnet that the bastion connects to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Bastion.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that bastion was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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
				Transform:   transform.From(bastionTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// Standard OCI columns
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
		},
	}
}

//// LIST FUNCTION

func listBastions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listBastions", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := bastionService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := bastion.ListBastionsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["id"] != nil {
		bastionId := equalQuals["id"].GetStringValue()
		request.BastionId = types.String(bastionId)
	}

	if equalQuals["lifecycle_state"] != nil {
		lifecycleState := equalQuals["lifecycle_state"].GetStringValue()
		request.BastionLifecycleState = bastion.ListBastionsBastionLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BastionClient.ListBastions(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, bastion := range response.Items {
			d.StreamListItem(ctx, bastion)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func getBastion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getBastion", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(bastion.BastionSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := bastionService(ctx, d, region)
	if err != nil {
		logger.Error("getBastion", "error_BastionService", err)
		return nil, err
	}

	request := bastion.GetBastionRequest{
		BastionId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BastionClient.GetBastion(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Bastion, nil
}

//// TRANSFORM FUNCTION

func bastionTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}

	switch d.HydrateItem.(type) {
	case bastion.Bastion:
		bastion := d.HydrateItem.(bastion.Bastion)
		freeformTags = bastion.FreeformTags
		definedTags = bastion.DefinedTags
	case bastion.BastionSummary:
		bastion := d.HydrateItem.(bastion.BastionSummary)
		freeformTags = bastion.FreeformTags
		definedTags = bastion.DefinedTags
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

	return tags, nil
}
