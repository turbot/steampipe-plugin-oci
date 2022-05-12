package oci

import (
	"context"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/v44/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v44/common"
	"github.com/oracle/oci-go-sdk/v44/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/connection"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

const matrixKeyRegion = "region"
const matrixKeyCompartment = "compartment"
const matrixKeyZone = "zone"

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

	if ociConfig.Regions != nil {
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

// BuildCompartmentList :: return a list of matrix items, one per compartment specified in the connection config
func BuildCompartmentList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {
	// cache compartment matrix
	cacheKey := "CompartmentList"

	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the compartments in the tenant
	compartments, err := listAllCompartments(ctx, pluginQueryData, connection)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	// validate compartment list
	matrix := make([]map[string]interface{}, len(compartments))
	for i, compartment := range compartments {
		matrix[i] = map[string]interface{}{matrixKeyCompartment: *compartment.Id}
	}
	// set CompartmentList cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// BuildCompartmentRegionList :: return a list of matrix items, one per region-compartment specified in the connection config
func BuildCompartementRegionList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {

	// cache compartment region matrix
	cacheKey := "CompartmentRegionList"

	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the compartments in the tenant
	compartments, err := listAllCompartments(ctx, pluginQueryData, connection)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	// retrieve regions from connection config
	ociConfig := GetConfig(connection)

	if ociConfig.Regions != nil {
		regions := GetConfig(connection).Regions

		if len(getInvalidRegions(regions)) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(getInvalidRegions(regions), ",") + ". Edit your connection configuration file and then restart Steampipe")
		}

		// validate regions list
		matrix := make([]map[string]interface{}, len(regions)*len(compartments))
		for i, region := range regions {
			for j, compartment := range compartments {
				matrix[len(compartments)*i+j] = map[string]interface{}{
					matrixKeyRegion:      region,
					matrixKeyCompartment: *compartment.Id,
				}
				plugin.Logger(ctx).Debug("listAllCompartments Matrix", (len(compartments)*i)+j, matrix[len(compartments)*i+j])
			}
		}

		// set CompartmentRegionList cache
		pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)
		return matrix
	}

	defaultMatrix := make([]map[string]interface{}, len(compartments))
	for j, compartment := range compartments {
		defaultMatrix[j] = map[string]interface{}{
			matrixKeyRegion:      getRegionFromEnvVar(),
			matrixKeyCompartment: *compartment.Id,
		}
		plugin.Logger(ctx).Debug("listAllCompartments MATRIX", j, defaultMatrix[j])
	}

	// set CompartmentRegionList cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, defaultMatrix)

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
		"sa-vinhedo-1",
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
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListCompartments(ctx, request)
		if err != nil {
			return nil, err
		}

		for _, compartment := range response.Items {
			if !helpers.StringSliceContains([]string{"CREATING", "DELETING", "DELETED"}, types.ToString(compartment.LifecycleState)) {
				compartments = append(compartments, compartment)
			}
		}

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

type zoneInfo struct {
	identity.AvailabilityDomain
	Region string
}

func listAllzones(ctx context.Context, d *plugin.QueryData, connection *plugin.Connection) ([]zoneInfo, error) {

	zonesList := []zoneInfo{}

	regions := GetConfig(connection).Regions

	if regions != nil {
		for _, region := range regions {
			session, err := identityServiceRegional(ctx, pluginQueryData, region)
			if err != nil {
				return nil, err
			}

			// The OCID of the tenancy containing the compartment.
			request := identity.ListAvailabilityDomainsRequest{
				CompartmentId: &session.TenancyID,
				RequestMetadata: oci_common.RequestMetadata{
					RetryPolicy: getDefaultRetryPolicy(d.Connection),
				},
			}

			response, err := session.IdentityClient.ListAvailabilityDomains(ctx, request)
			if err != nil {
				return nil, err
			}

			for _, zones := range response.Items {
				zonesList = append(zonesList, zoneInfo{zones, region})
			}
		}
		return zonesList, nil
	}
	region := getRegionFromEnvVar()
	session, err := identityServiceRegional(ctx, pluginQueryData, region)
	if err != nil {
		return nil, err
	}

	// The OCID of the tenancy containing the compartment.
	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: &session.TenancyID,
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.IdentityClient.ListAvailabilityDomains(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, zones := range response.Items {
		zonesList = append(zonesList, zoneInfo{zones, region})
	}
	return zonesList, nil
}

// BuildCompartmentZonalList :: return a list of matrix items, one per zone-compartment specified in the connection config
func BuildCompartementZonalList(ctx context.Context, connection *plugin.Connection) []map[string]interface{} {
	cacheKey := "CompartmentZonalList"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	compartments, err := listAllCompartments(ctx, pluginQueryData, connection)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	plugin.Logger(ctx).Debug("compartments", "compartments", compartments)

	zones, err := listAllzones(ctx, pluginQueryData, connection)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	matrix := make([]map[string]interface{}, len(zones)*len(compartments))

	for i, zone := range zones {
		for j, compartment := range compartments {
			matrix[len(compartments)*i+j] = map[string]interface{}{
				matrixKeyZone:        *zone.Name,
				matrixKeyCompartment: *compartment.Id,
				matrixKeyRegion:      zone.Region,
			}
			plugin.Logger(ctx).Debug("listAllCompartments Matrix", (len(compartments)*i)+j, matrix[len(compartments)*i+j])
		}
	}

	// set CompartmentZonalList cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// func getRegionFromEnvVar() (string, error) {
func getRegionFromEnvVar() string {
	if region, ok := os.LookupEnv("OCI_REGION"); ok {
		return region
	} else if region, ok = os.LookupEnv("OCI_CLI_REGION"); ok {
		return region
	}

	return getEnvSettingWithBlankDefault("region")
}

func getCloudGuardConfiguration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getCloudGuardConfiguration"
	if cachedData, ok := pluginQueryData.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(cloudguard.Configuration), nil
	}

	// Create Session
	session, err := cloudGuardService(ctx, d, "")
	if err != nil {
		return nil, err
	}

	request := cloudguard.GetConfigurationRequest{
		CompartmentId: types.String(session.TenancyID),
		RequestMetadata: oci_common.RequestMetadata{
			RetryPolicy: getDefaultRetryPolicy(d.Connection),
		},
	}

	response, err := session.CloudGuardClient.GetConfiguration(ctx, request)
	if err != nil {
		return nil, err
	}

	// set response cache
	pluginQueryData.ConnectionManager.Cache.Set(cacheKey, response.Configuration)

	return response.Configuration, nil
}

// Get home region
func getHomeRegion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	request := identity.ListRegionSubscriptionsRequest{
		TenancyId: &session.TenancyID,
	}

	// List all the subscribed regions for the tenant
	subscribedRegions, err := session.IdentityClient.ListRegionSubscriptions(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, subscribedRegion := range subscribedRegions.Items {
		if *subscribedRegion.IsHomeRegion {
			return *subscribedRegion.RegionName, nil
		}
	}
	return nil, nil
}
