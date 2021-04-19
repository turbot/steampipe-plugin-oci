package oci

import (
	"context"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v36/audit"
	"github.com/oracle/oci-go-sdk/v36/common"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAuditConfiguration(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_audit_configuration",
		Description: "OCI Audit Configuration",
		List: &plugin.ListConfig{
			Hydrate: listAuditConfigurations,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "event_id",
				Description: "The GUID of the event.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "event_type",
				Description: "The type of event that happened.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "source",
				Description: "The source of the event.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "retention_period_days",
				Description: "The retention period setting, specified in days. The minimum is 90, the maximum 365.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getAuditConfiguration,
				Transform:   transform.FromField("Configuration.RetentionPeriodDays"),
			},
			{
				Name:        "event_time",
				Description: "The time the event occurred.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("EventTime.Time"),
			},
			{
				Name:        "cloud_events_version",
				Description: "The version of the CloudEvents specification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "event_type_version",
				Description: "The version of the event type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content_type",
				Description: "The content type of the data contained in `data`.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "data",
				Description: "",
				Type:        proto.ColumnType_JSON,
			},

			// Standard OCI columns
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

func listAuditConfigurations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("listAuditConfigurations", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := auditService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	deafultTime, err := time.Parse(time.RFC3339, "2017-01-15T11:30:00Z")
	if err != nil {
		return nil, err
	}
	var deafultStartTime *common.SDKTime
	deafultStartTime.Time = deafultTime
	request := audit.ListEventsRequest{
		CompartmentId: types.String(compartment),
		StartTime:     deafultStartTime,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.AuditClient.(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, events := range response.Items {
			d.StreamListItem(ctx, events)
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

func getAuditConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAuditConfiguration")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Debug("getAuditConfiguration", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	// Create Session
	session, err := auditService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := audit.GetConfigurationRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.AuditClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Configuration, nil
}
