package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreBootVolumeAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_boot_volume_attachment",
		Description: "OCI Core Boot Volume Attachment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreBootVolumeAttachment,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreBootVolumeAttachments,
		},
		GetMatrixItem: BuildCompartementZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the boot volume attachment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name. Does not have to be unique, and it cannot be changed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "boot_volume_id",
				Description: "The OCID of the boot volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "instance_id",
				Description: "The OCID of the instance the boot volume is attached to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the boot volume attachment.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of an instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the boot volume was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "Encryption_In_Transit_Type",
				Description: "Refer the top-level definition of encryptionInTransitType.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_pv_encryption_in_transit_enabled",
				Description: "Whether in-transit encryption for the boot volume's paravirtualized attachment is enabled or not.",
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

func listCoreBootVolumeAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreBootVolumeAttachments", "Compartment", compartment, "OCI_Zone", zone)

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListBootVolumeAttachmentsRequest{
		CompartmentId:      types.String(compartment),
		AvailabilityDomain: types.String(zone),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListBootVolumeAttachments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeAttachment := range response.Items {
			d.StreamListItem(ctx, volumeAttachment)
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

func getCoreBootVolumeAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreBootVolumeAttachment")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreBootVolumeAttachment", "Compartment", compartment, "OCI_Zone", zone)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty volume attachment id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetBootVolumeAttachmentRequest{
		BootVolumeAttachmentId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ComputeClient.GetBootVolumeAttachment(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.BootVolumeAttachment, nil
}
