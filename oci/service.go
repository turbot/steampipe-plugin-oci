package oci

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	oci_common_auth "github.com/oracle/oci-go-sdk/v36/common/auth"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type session struct {
	TenancyID      string
	IdentityClient identity.IdentityClient
	ComputeClient  core.ComputeClient
	BlockstorageClient core.BlockstorageClient
}

// identityService returns the service client for OCI Identity service
func identityService(ctx context.Context, d *plugin.QueryData) (*session, error) {
	// if region == "" {
	// 	return nil, fmt.Errorf("region must be passed ACMService")
	// }
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("Identity-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		return nil, err
	}

	// provider := oci_common.CustomProfileConfigProvider(*ociConfig.ConfigPath, *ociConfig.Profile)
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:      tenantId,
		IdentityClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// coreBlockStorageService returns the service client for OCI Core BlockStorage Service
func coreBlockStorageService(ctx context.Context, d *plugin.QueryData) (*session, error) {

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("BlockstoragE-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider := oci_common.CustomProfileConfigProvider(*ociConfig.ConfigPath, *ociConfig.Profile)
	client, err := core.NewBlockstorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:     tenantID,
		BlockstorageClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

func coreComputeServiceRegional(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	// if region == "" {
	// 	return nil, fmt.Errorf("region must be passed ACMService")
	// }
	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ComputeRegional-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("coreComputeServiceRegional", "getProvider.Error", err)
		return nil, err
	}

	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:     tenantId,
		ComputeClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

func getProvider(ctx context.Context, d *connection.Manager, region string, config ociConfig) (oci_common.ConfigurationProvider, error) {

	cacheKey := "getProvider"
	// if provider is already cached, return it
	if cachedData, ok := d.Cache.Get(cacheKey); ok {
		return cachedData.(oci_common.ConfigurationProvider), nil
	}

	if region == "" && config.Regions != nil && len(config.Regions) > 0 {
		region = config.Regions[0]
	}

	if region == "" {
		region = getRegionFromEnvVar()
	}

	authType := "ApiKey"
	if config.Auth != nil && (*config.Auth != "ApiKey" && *config.Auth != "") {
		authType = *config.Auth
	}

	if authType == "SecurityToken" {
		return getProviderForSecurityToken(region, config)
	}

	if authType == "InstancePrincipal" {
		return getProviderForInstancePrincipal(region)
	}

	if authType == "ApiKey" {
		return getProviderForAPIkey(region, config)
	}

	regionInfo := oci_common.NewRawConfigurationProvider("", "", region, "", "", nil)
	provider, err := oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{regionInfo, oci_common.DefaultConfigProvider()})
	if err != nil {
		return nil, err
	}

	// save provider in cache
	d.Cache.Set(cacheKey, provider)

	return provider, nil
}

// https://github.com/oracle/oci-go-sdk/blob/master/example/helpers/helper.go#L127
func getDefaultRetryPolicy() *oci_common.RetryPolicy {
	// how many times to do the retry
	attempts := uint(5)

	// 429	TooManyRequests	You have issued too many requests to the
	// Oracle Cloud Infrastructure APIs in too short of an amount of time.	Yes, with backoff.

	// 500	InternalServerError	An internal server error occurred.	Yes, with backoff.

	// 503	ServiceUnavailable	The service is currently unavailable.	Yes, with backoff.
	// https: //docs.oracle.com/en-us/iaas/Content/API/References/apierrors.htm
	retryOnResponseCodes := func(r oci_common.OCIOperationResponse) bool {
		if r.Response.HTTPResponse() != nil {
			statusCode := strconv.Itoa(r.Response.HTTPResponse().StatusCode)
			return (r.Error != nil && helpers.StringSliceContains([]string{"429", "500", "503"}, statusCode))
		}
		return false
	}
	return getExponentialBackoffRetryPolicy(attempts, retryOnResponseCodes)
}

func getExponentialBackoffRetryPolicy(n uint, fn func(r oci_common.OCIOperationResponse) bool) *oci_common.RetryPolicy {
	// the duration between each retry operation, you might want to waite longer each time the retry fails
	exponentialBackoff := func(r oci_common.OCIOperationResponse) time.Duration {
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}
	policy := oci_common.NewRetryPolicy(n, fn, exponentialBackoff)
	return &policy
}

// connection "oci" {
//   tenancy_ocid = var.tenancy_ocid
//   config_file_profile= var.config_file_profile
// }

func getProviderForAPIkey(region string, config ociConfig) (oci_common.ConfigurationProvider, error) {

	// Check if all the attributes are available for API Key Authentication
	// connection "oci" {
	//   config_file_profile= var.config_file_profile
	// }

	regionInfo := oci_common.NewRawConfigurationProvider("", "", region, "", "", nil)

	if config.Profile != nil {
		configPath := ""
		if config.ConfigPath != nil {
			configPath = *config.ConfigPath
		}
		configProvider := oci_common.CustomProfileConfigProvider(configPath, *config.Profile)
		configProviderEnvironmentVariables := oci_common.ConfigurationProviderEnvironmentVariables("OCI_", "")

		return oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{regionInfo, configProvider, configProviderEnvironmentVariables})
	}

	// # Configure the Oracle Cloud Infrastructure provider with an API Key
	// connection "oci" {
	// 	tenancy_ocid = "tenancy_ocid"
	// 	user_ocid = "user_ocid"
	// 	fingerprint = "fingerprint"
	// 	private_key_path = "private_key_path"
	// 	regions = ["ap-mumbai-1", "us-ashburn-1"]
	// }

	if config.UserOCID != nil {
		pemFilePassword := ""
		pemFileContent := ""
		if config.PrivateKey != nil {
			pemFileContent = *config.PrivateKey
		}
		if config.PrivateKeyPath != nil {
			resolvedPath := expandPath(*config.PrivateKeyPath)
			pemFileData, err := ioutil.ReadFile(resolvedPath)
			if err != nil {
				return nil, fmt.Errorf("can not read private key from: '%s', Error: %q", *config.PrivateKeyPath, err)
			}
			pemFileContent = string(pemFileData)
		}

		if config.PrivateKeyPassword != nil {
			pemFilePassword = *config.PrivateKeyPassword
		}

		configProvider := oci_common.NewRawConfigurationProvider(*config.TenancyOCID, *config.UserOCID, region, *config.Fingerprint, pemFileContent, &pemFilePassword)
		configProviderEnvironmentVariables := oci_common.ConfigurationProviderEnvironmentVariables("OCI_", pemFilePassword)

		return oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{regionInfo, configProvider, configProviderEnvironmentVariables})
	}

	var providers []oci_common.ConfigurationProvider
	providers = append(providers, oci_common.DefaultConfigProvider())
	cliProvider, _ := getProviderFromCLIEnvironmentVariables()
	if cliProvider != nil {
		providers = append(providers, cliProvider)
	}

	// return default config in case connection config does not contain anything
	return oci_common.ComposingConfigurationProvider(providers)
}

// Check if all the attributes are available for SecurityToken Authentication
// connection "oci" {
//   auth = "SecurityToken"
//   config_file_profile= "config_file_profile"
// }
func getProviderForSecurityToken(region string, config ociConfig) (oci_common.ConfigurationProvider, error) {
	regionInfo := oci_common.NewRawConfigurationProvider("", "", region, "", "", nil)

	if config.Profile == nil {
		return nil, fmt.Errorf("\n\n'config_file_profile'must be set in the connection configuration for 'SecurityToken' authentication. Edit your connection configuration file and then restart Steampipe")
	}

	profileString := *config.Profile
	defaultPath := path.Join(getHomeFolder(), ".oci", "config")
	if err := checkProfile(profileString, defaultPath); err != nil {
		return nil, err
	}

	securityTokenBasedAuthConfigProvider := oci_common.CustomProfileConfigProvider(defaultPath, profileString)

	keyId, err := securityTokenBasedAuthConfigProvider.KeyID()
	if err != nil || !strings.HasPrefix(keyId, "ST$") {
		return nil, fmt.Errorf("security token is invalid")
	}

	return oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{regionInfo, securityTokenBasedAuthConfigProvider})
}

// # Configure the Oracle Cloud Infrastructure provider to use Instance Principal based authentication
// connection "oci" {
//   plugin 		= "oci"
//   auth 			= "InstancePrincipal"
//   region 		= [ "ap-mumbai-1" ]
// }

func getProviderForInstancePrincipal(region string) (oci_common.ConfigurationProvider, error) {

	// Used to modify InstancePrincipal auth clients so that `accept_local_certs` is honored for auth clients as well
	// These clients are created implicitly by SDK, and are not modified by the buildConfigureClientFn that usually does this for the other SDK clients
	instancePrincipalAuthClientModifier := func(client oci_common.HTTPRequestDispatcher) (oci_common.HTTPRequestDispatcher, error) {
		if acceptLocalCerts := getEnvSettingWithBlankDefault("accept_local_certs"); acceptLocalCerts != "" {
			if bool, err := strconv.ParseBool(acceptLocalCerts); err == nil {
				modifiedClient := buildHttpClient()
				modifiedClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = bool
				return modifiedClient, nil
			}
		}
		return client, nil
	}

	cfg, err := oci_common_auth.InstancePrincipalConfigurationForRegionWithCustomClient(oci_common.StringToRegion(region), instancePrincipalAuthClientModifier)
	if err != nil {
		return nil, err
	}

	return oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{cfg})
}

// cleans and expands the path if it contains a tilde , returns the expanded path or the input path as is if not expansion
// was performed
func expandPath(filepath string) string {
	if strings.HasPrefix(filepath, fmt.Sprintf("~%c", os.PathSeparator)) {
		filepath = path.Join(getHomeFolder(), filepath[2:])
	}
	return path.Clean(filepath)
}

func getHomeFolder() string {
	current, e := user.Current()
	if e != nil {
		//Give up and try to return something sensible
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return current.HomeDir
}

func checkProfile(profile string, path string) (err error) {
	var profileRegex = regexp.MustCompile(`^\[(.*)\]`)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	splitContent := strings.Split(content, "\n")
	for _, line := range splitContent {
		if match := profileRegex.FindStringSubmatch(line); len(match) > 1 && match[1] == profile {
			return nil
		}
	}

	return fmt.Errorf("configuration file did not contain profile: %s", profile)
}

func getEnvSettingWithBlankDefault(s string) string {
	return getEnvSettingWithDefault(s, "")
}

func getEnvSettingWithDefault(s string, dv string) string {
	v := os.Getenv("TF_VAR_" + s)
	if v != "" {
		return v
	}
	v = os.Getenv("OCI_" + s)
	if v != "" {
		return v
	}
	v = os.Getenv(s)
	if v != "" {
		return v
	}
	return dv
}

func getCLIEnvVariables(variableName string) string {
	v := os.Getenv("OCI_CLI_" + variableName)
	if v != "" {
		return v
	}
	v = os.Getenv("OCI_" + variableName)
	if v != "" {
		return v
	}
	return ""
}

func getProviderFromCLIEnvironmentVariables() (oci_common.ConfigurationProvider, error) {
	var providers []oci_common.ConfigurationProvider
	privateKeyPath := getCLIEnvVariables("KEY_FILE")
	pemFileContent := ""
	if privateKeyPath != "" {
		resolvedPath := expandPath(privateKeyPath)
		pemFileData, err := ioutil.ReadFile(resolvedPath)
		if err != nil {
			return nil, fmt.Errorf("can not read private key from: '%s', Error: %q", privateKeyPath, err)
		}
		pemFileContent = string(pemFileData)
	}

	cliApiKeyProvider := oci_common.NewRawConfigurationProvider(
		getCLIEnvVariables("TENANCY"),
		getCLIEnvVariables("USER"),
		getCLIEnvVariables("REGION"),
		getCLIEnvVariables("FINGERPRINT"),
		pemFileContent,
		types.String(""),
	)
	if cliApiKeyProvider != nil {
		providers = append(providers, cliApiKeyProvider)
	}

	cliFileWithProfileProvider, _ := oci_common.ConfigurationProviderFromFileWithProfile(
		getCLIEnvVariables("CONFIG_FILE"),
		getCLIEnvVariables("PROFILE"),
		getCLIEnvVariables(""),
	)

	if cliFileWithProfileProvider != nil {
		providers = append(providers, cliFileWithProfileProvider)
	}

	cliFromFileProvider, _ := oci_common.ConfigurationProviderFromFile(
		getCLIEnvVariables("CONFIG_FILE"),
		getCLIEnvVariables(""),
	)

	if cliFromFileProvider != nil {
		providers = append(providers, cliFromFileProvider)
	}

	if len(providers) > 0 {
		return oci_common.ComposingConfigurationProvider(providers)
	}
	return nil, nil
}

func buildHttpClient() (httpClient *http.Client) {
	httpClient = &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10000000000, // 10s
			}).DialContext,
			TLSHandshakeTimeout: 10000000000, // 10s
			TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
			Proxy:               http.ProxyFromEnvironment,
		},
	}
	return
}