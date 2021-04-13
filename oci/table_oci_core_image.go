package oci

import (
	"context"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableCoreImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "oci_core_image",
		Description: "OCI Core Image",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AnyColumn([]string{"id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404", "400"}),
			Hydrate:           getCoreImage,
		},
		List: &plugin.ListConfig{
			Hydrate: listCoreImages,
		},
		GetMatrixItem: BuildCompartementRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "display_name",
				Description: "A user-friendly name for the image. It does not have to be unique, and it's changeable.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The OCID of the image.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "lifecycle_state",
				Description: "The image's current state.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "time_created",
				Description: "The date and time the image was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("TimeCreated.Time"),
			},
			{
				Name:        "base_image_id",
				Description: "The OCID of the image originally used to launch the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromCamel(),
			},
			{
				Name:        "create_image_allowed",
				Description: "Indicates whether instances launched with this image can be used to create new images.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "launch_mode",
				Description: "Specifies the configuration mode for launching virtual machine (VM) instances.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "launch_options",
				Description: "LaunchOptions Options for tuning the compatibility and performance of VM shapes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "operating_system",
				Description: "The image's operating system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "operating_system_version",
				Description: "The image's operating system version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "size_in_mbs",
				Description: "The boot volume size for an instance launched from this image.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("SizeInMBs"),
			},
			{
				Name:        "agent_features",
				Description: "Oracle Cloud Agent features supported on the image.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.From(imageTags),
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
				Transform:   transform.FromField("CompartmentId"),
			},
		},
	}
}

//// LIST FUNCTION

func listCoreImages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("listCoreImages", "Compartment", compartment, "OCI_REGION", region)

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.ListImagesRequest{
		CompartmentId: types.String(compartment),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.ComputeClient.ListImages(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, image := range response.Items {
			d.StreamListItem(ctx, image)
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

func getCoreImage(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getImage")
	logger := plugin.Logger(ctx)
	region := plugin.GetMatrixItem(ctx)[matrixKeyRegion].(string)
	compartment := plugin.GetMatrixItem(ctx)[matrixKeyCompartment].(string)
	logger.Error("oci.getImage", "Compartment", compartment, "OCI_REGION", region)

	// Restrict the api call to only root compartment/ per region
	if !strings.HasPrefix(compartment, "ocid1.tenancy.oc1") {
		return nil, nil
	}

	id := d.KeyColumnQuals["id"].GetStringValue()

	// Create Session
	session, err := coreComputeService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	request := core.GetImageRequest{
		ImageId: types.String(id),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	response, err := session.ComputeClient.GetImage(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Image, nil
}

//// TRANSFORM FUNCTION

func imageTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	instance := d.HydrateItem.(core.Image)

	var tags map[string]interface{}

	if instance.FreeformTags != nil {
		tags = map[string]interface{}{}
		for k, v := range instance.FreeformTags {
			tags[k] = v
		}
	}

	if instance.DefinedTags != nil {
		if tags == nil {
			tags = map[string]interface{}{}
		}
		for _, v := range instance.DefinedTags {
			for key, value := range v {
				tags[key] = value
			}

		}
	}

	return tags, nil
}
