package oci

import (
	"context"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/bds"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// // TABLE DEFINITION
func tableBigDataServiceInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_big_data_service_instance",
		Description:      "OCI Big Data Service Instance",
		DefaultTransform: transform.FromCamel(),
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBigDataServiceInstance,
		},
		List: &plugin.ListConfig{
			Hydrate: listBigDataServiceInstances,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "compartment_id",
					Require: plugin.Optional,
				},
				{
					Name:    "lifecycle_state",
					Require: plugin.Optional,
				},
				{
					Name:    "display_name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The OCID of the Big Data Service resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The name of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cluster_profile",
				Description: "Profile of the Big Data Service cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that the Bds Instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "TimeUpdated",
				Description: "The time the cluster was updated, shown as an RFC 3339 formatted datetime string.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeUpdated.Time"),
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "lifecycle_state",
				Description: "The state of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_high_availability",
				Description: "Boolean flag specifying whether or not the cluster is highly available (HA).",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_secure",
				Description: "Boolean flag specifying whether or not the cluster should be set up as secure.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "is_cloud_sql_configured",
				Description: "Boolean flag specifying whether or not Cloud SQL should be configured.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "number_of_nodes",
				Description: "Number of nodes that form the cluster.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cluster_version",
				Description: "Version of the Hadoop distribution.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_by",
				Description: "The user who created the cluster.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "bootstrap_script_url",
				Description: "Pre-authenticated URL of the bootstrap script in Object Store that can be downloaded and executed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "kms_key_id",
				Description: "The OCID of the Key Management master encryption key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "network_config",
				Description: "Additional configuration of the user's network.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "cluster_details",
				Description: "Specific info about a Hadoop cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "cloud_sql_details",
				Description: "The information about added Cloud SQL capability.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "nodes",
				Description: "The list of nodes in the cluster.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigDataServiceInstance,
			},
			{
				Name:        "freeform_tags",
				Description: "Simple key-value pair that is applied without any predefined name, type, or scope.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "defined_tags",
				Description: "Defined tags for this resource. Each key is predefined and scoped to a namespace.",
				Type:        proto.ColumnType_JSON,
			},

			// Standard Steampipe columns
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bigDataServiceInstanceTags),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DisplayName"),
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
				Description: ColumnDescriptionTenantId,
				Type:        proto.ColumnType_STRING,
				Hydrate:     plugin.HydrateFunc(getTenantId).WithCache(),
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listBigDataServiceInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_big_data_service_instance.listBigDataServiceInstances", "Compartment", compartment, "OCI_REGION", region)

	equalQuals := d.EqualsQuals
	// Return nil, if given compartment_id doesn't match
	if equalQuals["compartment_id"] != nil && compartment != equalQuals["compartment_id"].GetStringValue() {
		return nil, nil
	}
	// Create Session
	session, err := bdsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_big_data_service_instance.listBigDataServiceInstances", "connection_error", err)
		return nil, err
	}

	//Build request parameters
	request := buildListBigDataServiceInstanceFilters(equalQuals)
	request.CompartmentId = types.String(compartment)
	request.Limit = types.Int(1000)
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BdsClient.ListBdsInstances(ctx, request)
		if err != nil {
			logger.Error("oci_big_data_service_instance.listBigDataServiceInstances", "api_error", err)
			return nil, err
		}
		for _, respItem := range response.Items {
			d.StreamListItem(ctx, respItem)

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

func getBigDataServiceInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := d.EqualsQualString(matrixKeyRegion)
	compartment := d.EqualsQualString(matrixKeyCompartment)
	logger.Debug("oci_big_data_service_instance.getBigDataServiceInstance", "Compartment", compartment, "OCI_REGION", region)
	if h.Item == nil && !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	var bdsInstanceId string
	if h.Item != nil {
		bdsInstanceId = *h.Item.(bds.BdsInstanceSummary).Id
	} else {
		bdsInstanceId = d.EqualsQualString("id")
	}

	if bdsInstanceId == "" {
		return nil, nil
	}

	request := bds.GetBdsInstanceRequest{
		BdsInstanceId: &bdsInstanceId,
	}

	// Create Session
	session, err := bdsService(ctx, d, region)
	if err != nil {
		logger.Error("oci_big_data_service_instance.getBigDataServiceInstance", "connection_error", err)
		return nil, err
	}
	request.RequestMetadata = common.RequestMetadata{
		RetryPolicy: getDefaultRetryPolicy(d.Connection),
	}

	response, err := session.BdsClient.GetBdsInstance(ctx, request)
	if err != nil {
		logger.Error("oci_big_data_service_instance.getBigDataServiceInstance", "api_error", err)
		return nil, err
	}
	return response.BdsInstance, nil
}

//// TRANSFORM FUNCTIONS

// Priority order for tags
// 1. Free-form tags
// 2. Defined Tags

func bigDataServiceInstanceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var freeformTags map[string]string
	var definedTags map[string]map[string]interface{}
	switch d.HydrateItem.(type) {
	case bds.BdsInstance:
		obj := d.HydrateItem.(bds.BdsInstance)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
	case bds.BdsInstanceSummary:
		obj := d.HydrateItem.(bds.BdsInstanceSummary)
		freeformTags = obj.FreeformTags
		definedTags = obj.DefinedTags
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

// Build additional list filters
func buildListBigDataServiceInstanceFilters(equalQuals plugin.KeyColumnEqualsQualMap) bds.ListBdsInstancesRequest {
	request := bds.ListBdsInstancesRequest{}

	if equalQuals["lifecycle_state"] != nil {
		request.LifecycleState = bds.BdsInstanceLifecycleStateEnum(equalQuals["lifecycle_state"].GetStringValue())
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	return request
}
