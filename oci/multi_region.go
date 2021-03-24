package oci

import (
	"context"
	"os"
	"strings"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const matrixKeyRegion = "region"
const matrixKeyCompartment = "compartment"

var pluginQueryData *plugin.QueryData

func init() {
	pluginQueryData = &plugin.QueryData{
		ConnectionManager: connection.NewManager(),
	}
}

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(_ context.Context, connection *plugin.Connection) []map[string]interface{} {
	// retrieve regions from connection config
	ociConfig := GetConfig(connection)

	if &ociConfig != nil && ociConfig.Regions != nil {
		regions := GetConfig(connection).Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(regions))
		for i, region := range regions {
			matrix[i] = map[string]interface{}{matrixKeyRegion: region}
		}
		return matrix
	}

	return []map[string]interface{}{
		{matrixKeyRegion: getRegionFromEnvVar()},
	}
}

// BuildCompartementRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildCompartementRegionList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {
	// get all the compartments in the tenant
	compartments, _ := listAllCompartments(ctx, pluginQueryData, connection)

	// retrieve regions from connection config
	ociConfig := GetConfig(connection)

	if &ociConfig != nil && ociConfig.Regions != nil {
		regions := GetConfig(connection).Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ","))
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(regions)*len(compartments))
		for i, region := range regions {
			for j, compartment := range compartments {
				matrix[len(compartments)*i+j] = map[string]interface{}{
					matrixKeyRegion:      region,
					matrixKeyCompartment: *compartment.Id,
				}
				plugin.Logger(ctx).Warn("MATRIX", (len(compartments)*i)+j, matrix[len(compartments)*i+j])
			}
		}
		return matrix
	}

	defaultMatrix := make([]map[string]interface{}, len(compartments))
	for j, compartment := range compartments {
		// plugin.Logger(ctx).Error("BuildCompartementRegionList", "compartment", compartment)
		defaultMatrix[j] = map[string]interface{}{
			matrixKeyRegion:      getRegionFromEnvVar(),
			matrixKeyCompartment: *compartment.Id,
		}
		plugin.Logger(ctx).Warn("MATRIX", j, defaultMatrix[j])
	}

	return defaultMatrix
}

func getInvalidRegions(regions []string) []string {
	ociRegions := []string{
		"ap-chiyoda-1",
		"ap-chuncheon-1",
		"ap-hyderabad-1",
		"ap-melbourne-1",
		"ap-mumbai-1",
		"ap-osaka-1",
		"ap-seoul-1",
		"ap-sydney-1",
		"ap-tokyo-1",
		"ca-montreal-1",
		"ca-toronto-1",
		"eu-amsterdam-1",
		"eu-frankfurt-1",
		"eu-zurich-1",
		"me-dubai-1",
		"me-jeddah-1",
		"sa-santiago-1",
		"sa-saopaulo-1",
		"uk-cardiff-1",
		"uk-gov-cardiff-1",
		"uk-gov-london-1",
		"uk-london-1",
		"us-ashburn-1",
		"us-gov-ashburn-1",
		"us-gov-chicago-1",
		"us-gov-phoenix-1",
		"us-langley-1",
		"us-luke-1",
		"us-phoenix-1",
		"us-sanjose-1",
	}

	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(ociRegions, region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
}

func listAllCompartments(ctx context.Context, d *plugin.QueryData, connection *plugin.Connection) ([]identity.Compartment, error) {
	// Create Session
	pluginQueryData.Connection = connection
	session, err := identityService(ctx, pluginQueryData)
	if err != nil {
		return nil, err
	}

	serviceCacheKey := "listAllCompartments"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.([]identity.Compartment), nil
	}

	// Add root tenant by default
	compartments := []identity.Compartment{
		{
			Id: &session.TenancyID,
		},
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListCompartmentsRequest{
		CompartmentId:          &session.TenancyID,
		CompartmentIdInSubtree: types.Bool(true),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListCompartments(ctx, request)
		if err != nil {
			return nil, err
		}
		compartments = append(compartments, response.Items...)

		if response.OpcNextPage != nil {
			request.Page = response.OpcNextPage
		} else {
			pagesLeft = false
		}
	}

	// save compartments in cache
	pluginQueryData.ConnectionManager.Cache.Set(serviceCacheKey, compartments)

	return compartments, err
}

// func getRegionFromEnvVar() (string, error) {
func getRegionFromEnvVar() string {
	region := os.Getenv("OCI_REGION")
	if region == "" {
		region = getEnvSettingWithBlankDefault("region")
	}
	return region
}
