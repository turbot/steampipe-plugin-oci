package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreVolumeAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_volume_attachment",
		Description: "OCI Core Volume Attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreVolumeAttachment,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreVolumeAttachments,
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
				{
					Name:    "volume_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the volume attachment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it cannot be changed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "volume_id",
				Description: "The OCID of the volume.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "instance_id",
				Description: "The OCID of the instance the volume is attached to.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the volume attachment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of an instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getCoreVolumeAttachmentFields,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "device",
				Description: "The device name.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "iscsi_login_state",
				Description: "The iscsi login state of the volume attachment.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "is_read_only",
				Description: "Whether the attachment was created in read-only mode.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "is_shareable",
				Description: "Whether the attachment should be created in shareable mode.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "is_pv_encryption_in_transit_enabled",
				Description: "Whether in-transit encryption for the data volume's paravirtualized attachment is enabled or not.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "is_multipath",
				Description: "Whether the attachment is multipath or not.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getCoreVolumeAttachmentFields,
			},

			// Standard Steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
				Transform:   transform.FromField("DisplayName"),
			},

			// Standard OCI columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
			},
			{
				Name:        "compartment_id",
				Description: ColumnDescriptionCompartment,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCoreVolumeAttachmentFields,
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

type volumeAttachmentInfo struct {
	// The availability domain of an instance.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string

	// The OCID of the compartment.
	CompartmentId *string

	// The OCID of the volume attachment.
	Id *string

	// The OCID of the instance the volume is attached to.
	InstanceId *string

	// The current state of the volume attachment.
	LifecycleState core.VolumeAttachmentLifecycleStateEnum

	// The date and time the volume was created, in the format defined by RFC3339 (https://tools.ietf.org/html/rfc3339).
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime

	// The OCID of the volume.
	VolumeId *string

	// The device name.
	Device *string

	// A user-friendly name. Does not have to be unique, and it cannot be changed.
	// Avoid entering confidential information.
	// Example: `My volume attachment`
	DisplayName *string

	// Whether the attachment was created in read-only mode.
	IsReadOnly *bool

	// Whether the attachment should be created in shareable mode. If an attachment
	// is created in shareable mode, then other instances can attach the same volume, provided
	// that they also create their attachments in shareable mode. Only certain volume types can
	// be attached in shareable mode. Defaults to false if not specified.
	IsShareable *bool

	// Whether in-transit encryption for the data volume's paravirtualized attachment is enabled or not.
	IsPvEncryptionInTransitEnabled *bool

	// Whether the attachment is multipath or not.
	IsMultipath *bool

	// The iscsi login state of the volume attachment. For a multipath volume attachment,
	// all iscsi sessions need to be all logged-in or logged-out to be in logged-in or logged-out state.
	IscsiLoginState core.VolumeAttachmentIscsiLoginStateEnum

	// Volume region
	Region string
}

//// LIST FUNCTION

func listCoreVolumeAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreVolumeAttachments", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := buildCoreVolumeAttachmentFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListVolumeAttachments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeAttachment := range response.Items {
			d.StreamListItem(ctx, volumeAttachment)

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

func getCoreVolumeAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreVolumeAttachment")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreVolumeAttachment", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty volume id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetVolumeAttachmentRequest{
		VolumeAttachmentId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ComputeClient.GetVolumeAttachment(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.VolumeAttachment, nil
}

func getCoreVolumeAttachmentFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(core.VolumeAttachment)

	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	plugin.Logger(ctx).Debug("getCoreVolumeAttachmentFields", "OCI_REGION", region)

	attachment := volumeAttachmentInfo{
		item.GetAvailabilityDomain(),
		item.GetCompartmentId(),
		item.GetId(),
		item.GetInstanceId(),
		item.GetLifecycleState(),
		item.GetTimeCreated(),
		item.GetVolumeId(),
		item.GetDevice(),
		item.GetDisplayName(),
		item.GetIsReadOnly(),
		item.GetIsShareable(),
		item.GetIsPvEncryptionInTransitEnabled(),
		item.GetIsMultipath(),
		item.GetIscsiLoginState(),
		region,
	}

	return attachment, nil
}

// Build additional filters
func buildCoreVolumeAttachmentFilters(equalQuals plugin.KeyColumnEqualsQualMap) core.ListVolumeAttachmentsRequest {
	request := core.ListVolumeAttachmentsRequest{}

	if equalQuals["availability_domain"] != nil {
		request.AvailabilityDomain = types.String(equalQuals["availability_domain"].GetStringValue())
	}
	if equalQuals["instance_id"] != nil {
		request.InstanceId = types.String(equalQuals["instance_id"].GetStringValue())
	}
	if equalQuals["volume_id"] != nil {
		request.VolumeId = types.String(equalQuals["volume_id"].GetStringValue())
	}

	return request
}
