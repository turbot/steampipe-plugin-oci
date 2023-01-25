package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreBlockVolumeReplica(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_block_volume_replica",
		Description: "OCI Core Block Volume Replica",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getCoreBlockVolumeReplica,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreBlockVolumeReplicas,
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
					Name:    "display_name",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartmentZonalList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The block volume replica's Oracle ID (OCID).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "display_name",
				Description: "A user-friendly name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of a block volume replica.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "block_volume_id",
				Description: "The OCID of the source block volume.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "availability_domain",
				Description: "The availability domain of the block volume replica.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the block volume replica was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},

			// other columns
			{
				Name:        "size_in_gbs",
				Description: "The size of the source block volume, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInGBs"),
			},
			{
				Name:        "time_last_synced",
				Description: "The date and time the block volume replica was last synced from the source block volume.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeLastSynced.Time"),
			},
			{
				Name:        "total_data_transferred_in_gbs",
				Description: "The total size of the data transferred from the source block volume to the block volume replica, in GBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("TotalDataTransferredInGBs"),
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
				Transform:   transform.From(volumeReplicaTags),
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

func listCoreBlockVolumeReplicas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listCoreBlockVolumeReplicas", "Compartment", compartment, "OCI_Zone", zone)

	equalQuals := d.KeyColumnQuals

	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}

	// Return nil, if given availability_domain doesn't match
	if equalQuals["availability_domain"] != nil && zone != equalQuals["availability_domain"].GetStringValue() {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListBlockVolumeReplicasRequest{
		AvailabilityDomain: types.String(zone),
		CompartmentId:      types.String(compartment),
		Limit:              types.Int(1000),
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
		request.LifecycleState = core.BlockVolumeReplicaLifecycleStateEnum(lifecycleState)
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BlockstorageClient.ListBlockVolumeReplicas(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, volumeReplica := range response.Items {
			d.StreamListItem(ctx, volumeReplica)

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

func getCoreBlockVolumeReplica(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	zone := plugin.GetMatrixItem(ctx)[matrixKeyZone].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getCoreBlockVolumeReplica", "Compartment", compartment, "OCI_Zone", zone)

	// Restrict the api call to only root compartment and one zone/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") || !strings.HasSuffix(zone, "AD-1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := coreBlockStorageService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetBlockVolumeReplicaRequest{
		BlockVolumeReplicaId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BlockstorageClient.GetBlockVolumeReplica(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.BlockVolumeReplica, nil
}

//// TRANSFORM FUNCTION

func volumeReplicaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	volumeReplica := d.HydrateItem.(core.BlockVolumeReplica)

	var tags map[string]interface{}

	if volumeReplica.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range volumeReplica.FreeformTags {
			tags[k] = v
		}
	}

	if volumeReplica.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range volumeReplica.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
