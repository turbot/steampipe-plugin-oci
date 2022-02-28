package oci

import (
	"context"
	"strconv"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVnicAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_vnic_attachment",
		Description: "OCI Core VNIC Attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getVnicAttachment,
		},
		List: &plugin.ListConfig{
			Hydrate:           listVnicAttachments,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "instance_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the VNIC attachment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name for the VNIC attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The OCID of the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the VNIC attachment. Possible values include: 'ATTACHING', 'ATTACHED', 'DETACHING', 'DETACHED'.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "private_ip",
				Description: "The private IP address of the primary `privateIp` object on the VNIC.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVnic,
				Transform:   transform.FromField("PrivateIp"),
			},
			{
				Name:        "public_ip",
				Description: "The public IP address of the VNIC, if one is assigned.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVnic,
				Transform:   transform.FromField("PublicIp"),
			},
			{
				Name:        "time_created",
				Description: "The date and time the VNIC attachment was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "hostname_label",
				Description: "The hostname for the VNIC's primary private IP.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVnic,
			},
			{
				Name:        "is_primary",
				Description: "Whether the VNIC is the primary VNIC (the VNIC that is automatically created and attached during instance launch).",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Hydrate:     getVnic,
			},
			{
				Name:        "mac_address",
				Description: "The MAC address of the VNIC.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVnic,
			},
			{
				Name:        "nic_index",
				Description: "The physical network interface card (NIC) the VNIC uses.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "skip_source_dest_check",
				Description: "Whether the source/destination check is disabled on the VNIC. Defaults to `false`, which means the check is performed.",
				Type:        proto.ColumnType_BOOL,
				Default:     false,
				Hydrate:     getVnic,
			},
			{
				Name:        "subnet_id",
				Description: "The OCID of the subnet to create the VNIC in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vlan_id",
				Description: "The OCID of the VLAN to create the VNIC in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vnic_id",
				Description: "The OCID of the VNIC.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "vnic_name",
				Description: "A user-friendly name for the VNIC.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getVnic,
				Transform:   transform.FromField("DisplayName"),
			},
			{
				Name:        "vlan_tag",
				Description: "The OCID of the VNIC.",
				Type:        proto.ColumnType_INT,
			},

			// JSON columns
			{
				Name:        "nsg_ids",
				Description: "A list of the OCIDs of the network security groups that the VNIC belongs to.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVnic,
			},

			// Tags
			{
				Name:        "defined_tags",
				Description: ColumnDescriptionDefinedTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVnic,
			},
			{
				Name:        "freeform_tags",
				Description: ColumnDescriptionFreefromTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVnic,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getVnic,
				Transform:   transform.From(vnicTags),
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
		},
	}
}

//// LIST FUNCTION

func listVnicAttachments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Trace("listVnicAttachments", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		logger.Error("listVnicAttachments", "compute_service_error", err)
		return nil, err
	}

	// Build request parameters
	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	request := core.ListVnicAttachmentsRequest{
		CompartmentId: types.String(compartment),
		Limit:         types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Check for additional filters
	if equalQuals["availability_domain"] != nil && equalQuals["availability_domain"].GetStringValue() != "" {
		availabilityDomain := equalQuals["availability_domain"].GetStringValue()
		request.AvailabilityDomain = types.String(availabilityDomain)
	}

	if equalQuals["instance_id"] != nil && equalQuals["instance_id"].GetStringValue() != "" {
		instanceId := equalQuals["instance_id"].GetStringValue()
		request.InstanceId = types.String(instanceId)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListVnicAttachments(ctx, request)
		if err != nil {
			logger.Error("listVnicAttachments", "list_vnic_attachments_error", err)
			return nil, err
		}

		for _, attachment := range response.Items {
			d.StreamListItem(ctx, attachment)

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

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVnicAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Trace("getVnicAttachment", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}
	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty VNIC attachment id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		logger.Error("getVnicAttachment", "compute_service_error", err)
		return nil, err
	}

	request := core.GetVnicAttachmentRequest{
		VnicAttachmentId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ComputeClient.GetVnicAttachment(ctx, request)
	if err != nil {
		logger.Error("getVnicAttachment", "get_vnic_attachment_error", err)
		return nil, err
	}

	return response.VnicAttachment, nil
}

func getVnic(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Trace("getVnic", "Compartment", compartment, "OCI_REGION", region)

	var vnicId string
	if h.Item != nil {
		vnicId = *h.Item.(core.VnicAttachment).VnicId
	} else {
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}

		vnicId = d.KeyColumnQuals["vnic_id"].GetStringValue()
	}

	// handle empty VNIC id in get call
	if vnicId == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreVirtualNetworkService(ctx, d, region)
	if err != nil {
		logger.Error("getVnic", "virtual_network_service_error", err)
		return nil, err
	}

	request := core.GetVnicRequest{
		VnicId: types.String(vnicId),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.VirtualNetworkClient.GetVnic(ctx, request)
	if err != nil {
		logger.Error("getVnic", "get_vnic_error", err)
		if ociErr, ok := err.(common.ServiceError); ok {
			if helpers.StringSliceContains([]string{"404"}, strconv.Itoa(ociErr.GetHTTPStatusCode())) {
				return nil, nil
			}
		}
		return nil, err
	}

	return response.Vnic, nil
}

//// TRANSFORM FUNCTION

// Priority order for tags
// 1. Defined Tags
// 2. Free-form tags
func vnicTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value != nil {
		vnic := d.HydrateItem.(core.Vnic)

		var tags map[string]interface{}

		if vnic.FreeformTags != nil {
			tags = map[string]interface{}{}
			for k, v := range vnic.FreeformTags {
				tags[k] = v
			}
		}

		if vnic.DefinedTags != nil {
			if tags == nil {
				tags = map[string]interface{}{}
			}
			for _, v := range vnic.DefinedTags {
				for key, value := range v {
					tags[key] = value
				}
			}
		}

		return tags, nil
	}

	return nil, nil
}
