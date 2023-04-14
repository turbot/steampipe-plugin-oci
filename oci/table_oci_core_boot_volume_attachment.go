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
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "availability_domain",
					Require: plugin.Optional,
				},
				{
					Name:    "boot_volume_id",
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
		GetMatrixItemFunc: BuildCompartementZonalList,
		Columns: commonColumnsForAllResource([]*plugin.Column{
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
				Name:        "encryption_in_transit_type",
				Description: "The type of the encryption in transit for the boot volume.",
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listCoreBootVolumeAttachments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	zone := d.EqualsQualString(matrixKeyZone)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("listCoreBootVolumeAttachments", "Compartment", compartment, "OCI_Zone", zone)

	equalQuals := d.EqualsQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Return nil, if given availability_domain doesn't match
	if equalQuals["availability_domain"] != nil && zone != equalQuals["availability_domain"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListBootVolumeAttachmentsRequest{
		AvailabilityDomain: types.String(zone),
		CompartmentId:      types.String(compartment),
		Limit:              types.Int(1000),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	// Check for additional filters
	if equalQuals["instance_id"] != nil {
		instanceId := equalQuals["instance_id"].GetStringValue()
		request.InstanceId = types.String(instanceId)
	}

	if equalQuals["boot_volume_id"] != nil {
		bootVolumeId := equalQuals["boot_volume_id"].GetStringValue()
		request.BootVolumeId = types.String(bootVolumeId)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListBootVolumeAttachments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeAttachment := range response.Items {
			d.StreamListItem(ctx, volumeAttachment)

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

func getCoreBootVolumeAttachment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCoreBootVolumeAttachment")
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	zone := d.EqualsQualString(matrixKeyZone)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("getCoreBootVolumeAttachment", "Compartment", compartment, "OCI_Zone", zone)

	// Restrict the api call to only root compartment and one zone/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetStringValue()

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
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.ComputeClient.GetBootVolumeAttachment(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.BootVolumeAttachment, nil
}
