package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/mysql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

//// TABLE DEFINITION

func tableMySQLHeatWaveCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_heat_wave_cluster",
		Description: "OCI MySQL heat wave cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"db_system_id"}),
			Hydrate:    getMySQLHeatWaveCluster,
		},
		List: &plugin.ListConfig{
			Hydrate:       listMySQLHeatWaveCluster,
			ParentHydrate: listMySQLDBSystems,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "db_system_id",
				Description: "The OCID of the parent DB System this HeatWave cluster is attached to.",
				Type:        proto.ColumnType_STRING,
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
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

func listMySQLHeatWaveCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	dbSystemId := h.Item.(mysql.DbSystemSummary).Id
	region := ociRegionNameFromId(*dbSystemId)
	logger.Debug("listMySQLHeatWaveCluster", "DB System", dbSystemId, "OCI_REGION", region)

	// Create Session
	session, err := mySQLDBSystemService(ctx, d, string(region))
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := mysql.GetHeatWaveClusterRequest{
		DbSystemId: dbSystemId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.MySQLDBSystemClient.GetHeatWaveCluster(ctx, request)
	if err != nil {
		return nil, err
	}

	heatWaveCluster := response.HeatWaveCluster
	d.StreamListItem(ctx, heatWaveCluster)

	return nil, nil
}

//// HYDRATE FUNCTION

func getMySQLHeatWaveCluster(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	dbSystemId := d.KeyColumnQuals["db_system_id"].GetStringValue()
	region := ociRegionNameFromId(dbSystemId)
	logger.Debug("getMySQLHeatWaveCluster", "DB System", dbSystemId, "OCI_REGION", region)

	// Create Session
	session, err := mySQLDBSystemService(ctx, d, string(region))
	if err != nil {
		return nil, err
	}

	// Build request parameters
	request := mysql.GetHeatWaveClusterRequest{
		DbSystemId: types.String(dbSystemId),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.MySQLDBSystemClient.GetHeatWaveCluster(ctx, request)
	if err != nil {
		return nil, err
	}

	heatWaveCluster := response.HeatWaveCluster

	return heatWaveCluster, nil
}
