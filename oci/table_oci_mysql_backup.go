package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/mysql"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableMySQLBackup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_mysql_backup",
		Description: "OCI MySQL Backup",
		List: &plugin.ListConfig{
			Hydrate: listMySQLBackups,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getMySQLBackup,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-supplied display name for the backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the backup.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "db_system_id",
				Description: "The OCID of the DB System the Backup is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the Backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The time the backup record was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMySQLBackup,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "backup_size_in_gbs",
				Description: "The size of the backup in GiBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("BackupSizeInGBs"),
			},
			{
				Name:        "backup_type",
				Description: "The type of backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "creation_type",
				Description: "If the backup was created automatically, or by a manual request.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A user-supplied description of the backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data_storage_size_in_gbs",
				Description: "Initial size of the data volume in GiBs.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DataStorageSizeInGBs"),
			},
			{
				Name:        "lifecycle_details",
				Description: "Additional information about the current lifecycleState.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getMySQLBackup,
			},
			{
				Name:        "mysql_version",
				Description: "The version of the DB System used for backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retention_in_days",
				Description: "Number of days to retain this backup.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "shape_name",
				Description: "The shape of the DB System instance used for backup.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_updated",
				Description: "The time at which the backup was updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getMySQLBackup,
				Transform:   transform.FromField("TimeUpdated.Time"),
			},
			{
				Name:        "db_system_snapshot",
				Description: "Snapshot of the DbSystem details at the time of the backup.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getMySQLBackup,
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

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(backupTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
			},

			// OCI standard columns
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
				Hydrate:     getMySQLBackup,
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

func listMySQLBackups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listMySQLBackups", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := mySQLDBBackupService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.ListBackupsRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.MySQLDBBackupClient.ListBackups(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, dbBackup := range response.Items {
			d.StreamListItem(ctx, dbBackup)
		}
		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMySQLBackup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getMySQLBackup", "Compartment", compartment, "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(mysql.BackupSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		// Restrict the api call to only root compartment/ per region
		if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
			return nil, nil
		}
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := mySQLDBBackupService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := mysql.GetBackupRequest{
		BackupId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.MySQLDBBackupClient.GetBackup(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Backup, nil
}

//// TRANSFORM FUNCTION

func backupTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	freeformTags := backupFreeformTags(d.HydrateItem)

	var tags map[string]interface{}

	if freeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range freeformTags {
			tags[k] = v
		}
	}

	definedTags := backupDefinedTags(d.HydrateItem)

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

func backupFreeformTags(item interface{}) map[string]string {
	switch item.(type) {
	case mysql.Backup:
		return item.(mysql.Backup).FreeformTags
	case mysql.BackupSummary:
		return item.(mysql.BackupSummary).FreeformTags
	}
	return nil
}

func backupDefinedTags(item interface{}) map[string]map[string]interface{} {
	switch item.(type) {
	case mysql.Backup:
		return item.(mysql.Backup).DefinedTags
	case mysql.BackupSummary:
		return item.(mysql.BackupSummary).DefinedTags
	}
	return nil
}
