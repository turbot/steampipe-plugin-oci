package oci

import (
	"context"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const matrixKeyRegion = "region"
const matrixKeyCompartment = "compartment"
const matrixKeyZone = "zone"

// BuildRegionList :: return a list of matrix items, one per region specified in the connection config
func BuildRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	// retrieve regions from connection config
	ociConfig := GetConfig(d.Connection)

	if ociConfig.Regions != nil {
		regions := GetConfig(d.Connection).Regions

		// fetch OCI regions
		validRegions, err := listOciAvailableRegions(ctx, d)
		if err != nil {
			panic(err)
		}
		invalidRegions := getInvalidRegions(regions, validRegions)
		if len(invalidRegions) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(invalidRegions, ","))
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
func BuildCompartmentList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	// cache compartment matrix
	cacheKey := "CompartmentList"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the compartments in the tenant
	compartments, err := listAllCompartments(ctx, d)
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
	d.ConnectionManager.Cache.Set(cacheKey, matrix)

	return matrix
}

// BuildCompartmentRegionList :: return a list of matrix items, one per region-compartment specified in the connection config
func BuildCompartementRegionList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// cache compartment region matrix
	cacheKey := "CompartmentRegionList"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	// get all the compartments in the tenant
	compartments, err := listAllCompartments(ctx, d)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	// retrieve regions from connection config
	ociConfig := GetConfig(d.Connection)

	if ociConfig.Regions != nil {
		regions := GetConfig(d.Connection).Regions

		// fetch OCI regions
		validRegions, err := listOciAvailableRegions(ctx, d)
		if err != nil {
			panic(err)
		}
		invalidRegions := getInvalidRegions(regions, validRegions)
		if len(invalidRegions) > 0 {
			panic("\n\nConnection config have invalid regions: " + strings.Join(invalidRegions, ","))
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
		d.ConnectionManager.Cache.Set(cacheKey, matrix)
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
	d.ConnectionManager.Cache.Set(cacheKey, defaultMatrix)

	return defaultMatrix
}

func getInvalidRegions(regions []string, ociRegions []string) []string {
	invalidRegions := []string{}
	for _, region := range regions {
		if !helpers.StringSliceContains(ociRegions, region) {
			invalidRegions = append(invalidRegions, region)
		}
	}
	return invalidRegions
}

func listAllCompartments(ctx context.Context, d *plugin.QueryData) ([]identity.Compartment, error) {
	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	serviceCacheKey := "listAllCompartments"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
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

	endpointRegion := GetConfig(d.Connection).Regions[0]
	pagesLeft := true
	for pagesLeft {
		response, err := session.IdentityClient.ListCompartments(ctx, request)
		if err != nil {
			if strings.Contains(err.Error(), "no such host") {
				panic("\n\nConnection config has invalid region: " + endpointRegion + ". Edit your connection configuration file and then restart Steampipe")
			}
			plugin.Logger(ctx).Error("listAllCompartments", "ListCompartments", err)
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
	d.ConnectionManager.Cache.Set(serviceCacheKey, compartments)

	return compartments, err
}

type zoneInfo struct {
	identity.AvailabilityDomain
	Region string
}

func listAllzones(ctx context.Context, d *plugin.QueryData) ([]zoneInfo, error) {

	zonesList := []zoneInfo{}

	regions := GetConfig(d.Connection).Regions

	if regions != nil {
		for _, region := range regions {
			session, err := identityServiceRegional(ctx, d, region)
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
	session, err := identityServiceRegional(ctx, d, region)
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
func BuildCompartementZonalList(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	cacheKey := "CompartmentZonalList"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	compartments, err := listAllCompartments(ctx, d)
	if err != nil {
		if strings.Contains(err.Error(), "proper configuration for region") || strings.Contains(err.Error(), "OCI_REGION") {
			panic("\n\n'regions' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
		}
		panic(err)
	}

	plugin.Logger(ctx).Debug("compartments", "compartments", compartments)

	zones, err := listAllzones(ctx, d)
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
	d.ConnectionManager.Cache.Set(cacheKey, matrix)

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
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
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
	d.ConnectionManager.Cache.Set(cacheKey, response.Configuration)

	return response.Configuration, nil
}

// List out the regions supported by Oracle Cloud
func listOciAvailableRegions(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	logger := plugin.Logger(ctx)

	cacheKey := "OciRegionList"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]string), nil
	}

	endpointRegion := GetConfig(d.Connection).Regions[0]

	// Create Session
	session, err := identityService(ctx, d)
	if err != nil {
		return nil, err
	}

	var regionNames []string

	regions, err := session.IdentityClient.ListRegions(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			panic("\n\nConnection config has invalid region: " + endpointRegion + ". Edit your connection configuration file and then restart Steampipe")
		}
		logger.Error("listOciAvailableRegions", "ListRegions", err)
		return nil, err
	}

	for _, region := range regions.Items {
		regionNames = append(regionNames, *region.Name)
	}

	d.ConnectionManager.Cache.Set(cacheKey, regionNames)

	return regionNames, nil
}
