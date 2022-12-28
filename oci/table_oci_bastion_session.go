package oci

import (
	"context"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableBastionSession(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "oci_bastion_session",
		Description:      "OCI Bastion Session",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			Hydrate: listBastionSessions,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "bastion_id",
					Require: plugin.Required,
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
		GetMatrixItemFunc: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The unique identifier (OCID) of the session.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bastion_id",
				Description: "The unique identifier (OCID) of the bastion that is hosting this session.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bastion_name",
				Description: "The name of the bastion that is hosting this session.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_resource_details",
				Description: "Details about a bastion session's target resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_details",
				Description: "Public key details for a bastion session.",
				Hydrate:     getBastionSession,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "session_ttl_in_seconds",
				Description: "The amount of time the session can remain active.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "display_name",
				Description: "The name of the session.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "bastion_user_name",
				Description: "The username that the session uses to connect to the target resource.",
				Hydrate:     getBastionSession,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssh_metadata",
				Description: "The connection message for the session.",
				Hydrate:     getBastionSession,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "key_type",
				Description: "The type of the key used to connect to the session. PUB is a standard public key in OpenSSH format.",
				Hydrate:     getBastionSession,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "bastion_public_host_key_info",
				Description: "The public key of the bastion host. You can use this to verify that you're connecting to the correct bastion.",
				Hydrate:     getBastionSession,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "lifecycle_state",
				Description: "The current state of the session.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "Time that bastion was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
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

func listBastionSessions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQuals["region"].GetStringValue()
	equalQuals := d.KeyColumnQuals
	bastionId := equalQuals["bastion_id"].GetStringValue()

	// handle empty keyId, endpoint and region in list call
	if bastionId == "" {
		return nil, nil
	}

	// Create Session
	session, err := bastionService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := bastion.ListSessionsRequest{
		BastionId: types.String(bastionId),
		Limit:     types.Int(100),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	if equalQuals["display_name"] != nil {
		request.DisplayName = types.String(equalQuals["display_name"].GetStringValue())
	}

	// Check for limit
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < int64(*request.Limit) {
			request.Limit = types.Int(int(*limit))
		}
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.BastionClient.ListSessions(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, bastionSession := range response.Items {
			d.StreamListItem(ctx, bastionSession)

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

// // HYDRATE FUNCTION
func getBastionSession(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	logger.Debug("getBastionSession", "OCI_REGION", region)

	var id string
	if h.Item != nil {
		id = *h.Item.(bastion.SessionSummary).Id
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()

	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	// Create Session
	session, err := bastionService(ctx, d, region)
	if err != nil {
		logger.Error("getBastionSession", "error_BastionService", err)
		return nil, err
	}

	request := bastion.GetSessionRequest{
		SessionId: types.String(id),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.BastionClient.GetSession(ctx, request)

	if err != nil {
		return nil, err
	}
	return response.Session, nil
}
