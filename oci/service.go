package oci

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/oracle/oci-go-sdk/v36/apigateway"
	"github.com/oracle/oci-go-sdk/v36/audit"
	"github.com/oracle/oci-go-sdk/v36/autoscaling"
	"github.com/oracle/oci-go-sdk/v36/cloudguard"
	oci_common "github.com/oracle/oci-go-sdk/v36/common"
	oci_common_auth "github.com/oracle/oci-go-sdk/v36/common/auth"
	"github.com/oracle/oci-go-sdk/v36/core"
	"github.com/oracle/oci-go-sdk/v36/dns"
	"github.com/oracle/oci-go-sdk/v36/events"
	"github.com/oracle/oci-go-sdk/v36/filestorage"
	"github.com/oracle/oci-go-sdk/v36/functions"
	"github.com/oracle/oci-go-sdk/v36/identity"
	"github.com/oracle/oci-go-sdk/v36/keymanagement"
	"github.com/oracle/oci-go-sdk/v36/logging"
	"github.com/oracle/oci-go-sdk/v36/mysql"
	"github.com/oracle/oci-go-sdk/v36/objectstorage"
	"github.com/oracle/oci-go-sdk/v36/ons"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type session struct {
	TenancyID                      string
	ApiGatewayClient               apigateway.ApiGatewayClient
	AuditClient                    audit.AuditClient
	AutoScalingClient              autoscaling.AutoScalingClient
	BlockstorageClient             core.BlockstorageClient
	CloudGuardClient               cloudguard.CloudGuardClient
	ComputeClient                  core.ComputeClient
	DnsClient                      dns.DnsClient
	EventsClient                   events.EventsClient
	FileStorageClient              filestorage.FileStorageClient
	FunctionsManagementClient      functions.FunctionsManagementClient
	IdentityClient                 identity.IdentityClient
	KmsManagementClient            keymanagement.KmsManagementClient
	KmsVaultClient                 keymanagement.KmsVaultClient
	LoggingManagementClient        logging.LoggingManagementClient
	MysqlClient                    mysql.DbSystemClient
	NotificationControlPlaneClient ons.NotificationControlPlaneClient
	NotificationDataPlaneClient    ons.NotificationDataPlaneClient
	ObjectStorageClient            objectstorage.ObjectStorageClient
	VirtualNetworkClient           core.VirtualNetworkClient
}

// apiGatewayService returns the service client for OCI ApiGateway service
func apiGatewayService(ctx context.Context, d *plugin.QueryData) (*session, error) {

	serviceCacheKey := fmt.Sprintf("ApiGateway-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		return nil, err
	}

	// get apigateway service client
	client, err := apigateway.NewApiGatewayClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:        tenantId,
		ApiGatewayClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// auditService returns the service client for OCI Audit service
func auditService(ctx context.Context, d *plugin.QueryData) (*session, error) {

	serviceCacheKey := fmt.Sprintf("audit-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		return nil, err
	}

	// get audit service client
	client, err := audit.NewAuditClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:   tenantId,
		AuditClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// autoScalingService returns the service client for OCI Auto Scaling Service
func autoScalingService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("AutoScaling-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("autoScalingService", "getProvider.Error", err)
		return nil, err
	}

	client, err := autoscaling.NewAutoScalingClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:         tenantId,
		AutoScalingClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// identityService returns the service client for OCI Identity service
func identityService(ctx context.Context, d *plugin.QueryData) (*session, error) {

	serviceCacheKey := fmt.Sprintf("Identity-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		return nil, err
	}

	// get identity service client
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
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

// identityServiceRegional returns the service client for OCI Logging Management Service
func identityServiceRegional(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("LoggingManagement-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("loggingManagementService", "getProvider.Error", err)
		return nil, err
	}

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

// loggingManagementService returns the service client for OCI Logging Management Service
func loggingManagementService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("LoggingManagement-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("loggingManagementService", "getProvider.Error", err)
		return nil, err
	}

	client, err := logging.NewLoggingManagementClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:               tenantId,
		LoggingManagementClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// coreBlockStorageService returns the service client for OCI Core BlockStorage Service
func coreBlockStorageService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("BlockStorage-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("coreBlockStorageService", "getProvider.Error", err)
		return nil, err
	}

	client, err := core.NewBlockstorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:          tenantId,
		BlockstorageClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// eventsService returns the service client for OCI Events Service
func eventsService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("Events-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("eventsService", "getProvider.Error", err)
		return nil, err
	}

	client, err := events.NewEventsClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:    tenantId,
		EventsClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// fileStorageService returns the service client for OCI File Storage Service
func fileStorageService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("FileStorage-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("fileStorageService", "getProvider.Error", err)
		return nil, err
	}
	client, err := filestorage.NewFileStorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:         tenantID,
		FileStorageClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// functionsManagementService returns the service client for OCI Core VirtualNetwork Service
func functionsManagementService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("FunctionsManagement-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("functionsManagementService", "getProvider.Error", err)
		return nil, err
	}

	client, err := functions.NewFunctionsManagementClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:                 tenantID,
		FunctionsManagementClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// kmsManagementService returns the service client for OCI KMS Management Service
func kmsManagementService(ctx context.Context, d *plugin.QueryData, region string, endpoint string) (*session, error) {
	logger := plugin.Logger(ctx)

	// Cache the connection at vault level
	serviceCacheKey := fmt.Sprintf("KeyManagement-%s-%s", region, endpoint)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}
	// get oci config info
	ociConfig := GetConfig(d.Connection)
	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("kmsManagementService", "getProvider.Error", err)
		return nil, err
	}

	client, err := keymanagement.NewKmsManagementClientWithConfigurationProvider(provider, endpoint)
	if err != nil {
		return nil, err
	}
	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}
	sess := &session{
		TenancyID:           tenantId,
		KmsManagementClient: client,
	}
	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)
	return sess, nil
}

// kmsVaultService returns the service client for OCI KMS Vault Service
func kmsVaultService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("Vault-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("kmsVaultService", "getProvider.Error", err)
		return nil, err
	}

	client, err := keymanagement.NewKmsVaultClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:      tenantId,
		KmsVaultClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// objectStorageService returns the service client for OCI Object Storage service
func objectStorageService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("ObjectStorage-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("objectStorageService", "getProvider.Error", err)
		return nil, err
	}

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:           tenantId,
		ObjectStorageClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// onsNotificationControlPlaneService returns the service client for OCI Notification Control Plane service
func onsNotificationControlPlaneService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("NotificationControlPlane-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("onsNotificationControlPlaneService", "getProvider.Error", err)
		return nil, err
	}

	// get notification service client
	client, err := ons.NewNotificationControlPlaneClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:                      tenantId,
		NotificationControlPlaneClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// onsNotificationDataPlaneService returns the service client for OCI Notification Data Plane service
func onsNotificationDataPlaneService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("NotificationDataPlane-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("notificationDataPlaneService", "getProvider.Error", err)
		return nil, err
	}

	// get notification data plane service client
	client, err := ons.NewNotificationDataPlaneClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
	tenantId, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:                   tenantId,
		NotificationDataPlaneClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// coreComputeService returns the service client for OCI Core Compute service
func coreComputeService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the service?
	serviceCacheKey := fmt.Sprintf("ComputeRegional-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info from steampipe connection
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("coreComputeServiceRegional", "getProvider.Error", err)
		return nil, err
	}

	// get compute service client
	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	// get tenant ocid from provider
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

// coreVirtualNetworkService returns the service client for OCI Core VirtualNetwork Service
func coreVirtualNetworkService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("VirtualNetwork-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("coreVirtualNetworkService", "getProvider.Error", err)
		return nil, err
	}

	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:            tenantID,
		VirtualNetworkClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// cloudGuardService returns the service client for OCI Cloud Guard Service
func cloudGuardService(ctx context.Context, d *plugin.QueryData) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("CloudGuard-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		logger.Error("cloudGuardService", "getProvider.Error", err)
		return nil, err
	}

	client, err := cloudguard.NewCloudGuardClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:        tenantID,
		CloudGuardClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// dnsService returns the service client for OCI DNS Service
func dnsService(ctx context.Context, d *plugin.QueryData) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("dns-%s", "region")
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, "", ociConfig)
	if err != nil {
		logger.Error("DNSService", "getProvider.Error", err)
		return nil, err
	}

	client, err := dns.NewDnsClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID: tenantID,
		DnsClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// mysqlService returns the service client for OCI Mysql Service
func mysqlService(ctx context.Context, d *plugin.QueryData, region string) (*session, error) {
	logger := plugin.Logger(ctx)
	serviceCacheKey := fmt.Sprintf("mysql-%s", region)
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*session), nil
	}

	// get oci config info
	ociConfig := GetConfig(d.Connection)

	provider, err := getProvider(ctx, d.ConnectionManager, region, ociConfig)
	if err != nil {
		logger.Error("mysqlService", "getProvider.Error", err)
		return nil, err
	}

	client, err := mysql.NewDbSystemClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	tenantID, err := provider.TenancyOCID()
	if err != nil {
		return nil, err
	}

	sess := &session{
		TenancyID:   tenantID,
		MysqlClient: client,
	}

	// save session in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, sess)

	return sess, nil
}

// get the configurtion provider for the OCI plugin connection to intract with API's
func getProvider(_ context.Context, d *connection.Manager, region string, config ociConfig) (oci_common.ConfigurationProvider, error) {

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

/*
	#  Configure the Oracle Cloud Infrastructure provider with an API Key / or a profile
	connection "oci" {
		config_file_profile = "DEFAULT"
		config_path = "~/Desktop/config"
		regions = ["ap-mumbai-1", "us-ashburn-1"]
	}

	connection "oci" {
		tenancy_ocid = "tenancy_ocid"
		user_ocid = "user_ocid"
		fingerprint = "fingerprint"
		private_key_path = "private_key_path"
		regions = ["ap-mumbai-1", "us-ashburn-1"]
	}
*/
func getProviderForAPIkey(region string, config ociConfig) (oci_common.ConfigurationProvider, error) {

	// config provider with region info
	regionInfo := oci_common.NewRawConfigurationProvider("", "", region, "", "", nil)

	if config.Profile != nil {
		configPath := ""
		if config.ConfigPath != nil {
			configPath = *config.ConfigPath
		}

		// If the ~/.steampipe/config/oci.spc contains a profile, gets provider for it
		configProvider := oci_common.CustomProfileConfigProvider(configPath, *config.Profile)
		configProviderEnvironmentVariables := oci_common.ConfigurationProviderEnvironmentVariables("OCI_", "")

		return oci_common.ComposingConfigurationProvider([]oci_common.ConfigurationProvider{regionInfo, configProvider, configProviderEnvironmentVariables})
	}

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
	providers = append(providers, regionInfo, oci_common.DefaultConfigProvider())
	cliProvider, _ := getProviderFromCLIEnvironmentVariables()
	if cliProvider != nil {
		providers = append(providers, cliProvider)
	}

	// return default config in case connection config does not contain anything
	return oci_common.ComposingConfigurationProvider(providers)
}

/*
	# Provider for SecurityToken Authentication
	connection "oci" {
		auth = "SecurityToken"
		config_file_profile= "config_file_profile"
	}
*/
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

/*
# Provider for Instance Principal based authentication
	connection "oci" {
		plugin 		= "oci"
		auth 			= "InstancePrincipal"
		region 		= [ "ap-mumbai-1" ]
	}
*/
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

// cleans and expands the path if it contains a tilde,
// returns the expanded path or the input path as is if not expansion was performed
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

// Check for the profile in config file
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

// Get the value of environment variables
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

// Get the value of environment variables of OCI CLI
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

// Get the provider from OCI CLI environment variables
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
