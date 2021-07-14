package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the volume attachment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it cannot be changed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "volume_id",
				Description: "The OCID of the volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "instance_id",
				Description: "The OCID of the instance the volume is attached to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the volume attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of an instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "device",
				Description: "The device name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_read_only",
				Description: "Whether the attachment was created in read-only mode.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_shareable",
				Description: "Whether the attachment should be created in shareable mode.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_pv_encryption_in_transit_enabled",
				Description: "Whether in-transit encryption for the data volume's paravirtualized attachment is enabled or not.",
				Type:        proto.ColumnType_BOOL,
			},

			// Standard Steampipe columns
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

	// Volume region
	Region string
}

//// LIST FUNCTION

func listCoreVolumeAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreVolumeAttachments", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListVolumeAttachmentsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListVolumeAttachments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeAttachment := range response.Items {
			d.StreamListItem(ctx, volumeAttachmentInfo{volumeAttachment.GetAvailabilityDomain(), volumeAttachment.GetCompartmentId(), volumeAttachment.GetId(), volumeAttachment.GetInstanceId(), volumeAttachment.GetLifecycleState(), volumeAttachment.GetTimeCreated(), volumeAttachment.GetVolumeId(), volumeAttachment.GetDevice(), volumeAttachment.GetDisplayName(), volumeAttachment.GetIsReadOnly(), volumeAttachment.GetIsShareable(), volumeAttachment.GetIsPvEncryptionInTransitEnabled(), region})
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

	volumeAttachment := response.VolumeAttachment

	return volumeAttachmentInfo{volumeAttachment.GetAvailabilityDomain(), volumeAttachment.GetCompartmentId(), volumeAttachment.GetId(), volumeAttachment.GetInstanceId(), volumeAttachment.GetLifecycleState(), volumeAttachment.GetTimeCreated(), volumeAttachment.GetVolumeId(), volumeAttachment.GetDevice(), volumeAttachment.GetDisplayName(), volumeAttachment.GetIsReadOnly(), volumeAttachment.GetIsShareable(), volumeAttachment.GetIsPvEncryptionInTransitEnabled(), region}, nil
}
