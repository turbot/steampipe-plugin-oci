package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableOciMySQLHeatWaveCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_heat_wave_cluster",
		Description: "OCI MySQL heat wave cluster",
		List: &plugin.ListConfig{
			Hydrate:       listMySQLHeatWaveCluster,
			ParentHydrate: listMySQLDBSystems,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "db_system_id",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "db_system_id",
				Description: "The OCID of the parent DB System this HeatWave cluster is attached to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "time_created",
				Description: "The date and time the HeatWave cluster was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "cluster_size",
				Description: "The number of analytics-processing compute instances, of the specified shape, in the HeatWave cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycleState.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the HeatWave cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "shape_name",
				Description: "The shape determines resources to allocate to the HeatWave nodes - CPU cores, memory.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The time the HeatWave cluster was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "cluster_nodes",
				Description: "A HeatWave node is a compute host that is part of a HeatWave cluster.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

//// LIST FUNCTION

func listMySQLHeatWaveCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	dbSystem := h.Item.(mysql.DbSystemSummary)
	region := ociRegionNameFromId(*dbSystem.Id)
	logger.Debug("listMySQLHeatWaveCluster", "DB System", dbSystem, "OCI_REGION", region)

	if d.KeyColumnQualString("db_system_id") != "" && d.KeyColumnQualString("db_system_id") != *dbSystem.Id {
		return nil, nil
	}

	// Create Session
	session, err := mySQLDBSystemService(ctx, d, string(region))
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := mysql.GetHeatWaveClusterRequest{
		DbSystemId: dbSystem.Id,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	// Ignore if there is no heat wave cluster associated to DB System
	if !*dbSystem.IsHeatWaveClusterAttached {
		return nil, nil
	}

	response, err := session.MySQLDBSystemClient.GetHeatWaveCluster(ctx, request)
	if err != nil {
		return nil, nil
	}
	d.StreamListItem(ctx, response.HeatWaveCluster)

	return nil, nil
}
