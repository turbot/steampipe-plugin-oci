/*
Package oci implements a steampipe plugin for OCI.

This plugin provides data that Steampipe uses to present foreign
tables that represent Oracle Cloud Infrastructure resources.
*/
package oci

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-oci"

// Plugin creates this (oci) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"oci_autoscaling_auto_scaling_configuration": tableAutoScalingConfiguration(ctx),
			"oci_cloud_guard_configuration":              tableCloudGuardConfiguration(ctx),
			"oci_cloud_guard_detector_recipe":            tableCloudGuardDetectorRecipe(ctx),
			"oci_cloud_guard_managed_list":               tableCloudGuardManagedList(ctx),
			"oci_cloud_guard_responder_recipe":           tableCloudGuardResponderRecipe(ctx),
			"oci_cloud_guard_target":                     tableCloudGuardTarget(ctx),
			"oci_core_dhcp_options":                      tableCoreDhcpOptions(ctx),
			"oci_core_drg":                               tableCoreDrg(ctx),
			"oci_core_image":                             tableCoreImage(ctx),
			"oci_core_instance":                          tableCoreInstance(ctx),
			"oci_core_internet_gateway":                  tableCoreInternetGateway(ctx),
			"oci_core_local_peering_gateway":             tableCoreLocalPeeringGateway(ctx),
			"oci_core_nat_gateway":                       tableCoreNatGateway(ctx),
			"oci_core_network_security_group":            tableCoreNetworkSecurityGroup(ctx),
			"oci_core_route_table":                       tableCoreRouteTable(ctx),
			"oci_core_security_list":                     tableCoreSecurityList(ctx),
			"oci_core_service_gateway":                   tableCoreServiceGateway(ctx),
			"oci_core_subnet":                            tableCoreSubnet(ctx),
			"oci_core_vcn":                               tableCoreVcn(ctx),
			"oci_core_volume":                            tableCoreVolume(ctx),
			"oci_core_volume_backup":                     tableCoreVolumeBackup(ctx),
			"oci_core_volume_backup_policy":              tableCoreVolumeBackupPolicy(ctx),
			"oci_dns_tsig_key":                           tableDnsTsigKey(ctx),
			"oci_dns_zone":                               tableDnsZone(ctx),
			"oci_events_rule":                            tableEventsRule(ctx),
			"oci_file_storage_file_system":               tableFileStorageFileSystem(ctx),
			"oci_functions_application":                  tableFunctionsApplication(ctx),
			"oci_identity_auth_token":                    tableIdentityAuthToken(ctx),
			"oci_identity_authentication_policy":         tableIdentityAuthenticationPolicy(ctx),
			"oci_identity_compartment":                   tableIdentityCompartment(ctx),
			"oci_identity_customer_secret_key":           tableIdentityCustomerSecretKey(ctx),
			"oci_identity_dynamic_group":                 tableIdentityDynamicGroup(ctx),
			"oci_identity_group":                         tableIdentityGroup(ctx),
			"oci_identity_network_source":                tableIdentityNetworkSource(ctx),
			"oci_identity_policy":                        tableIdentityPolicy(ctx),
			"oci_identity_tag_namespace":                 tableIdentityTagNamespace(ctx),
			"oci_identity_tenancy":                       tableIdentityTenancy(ctx),
			"oci_identity_user":                          tableIdentityUser(ctx),
			"oci_kms_vault":                              tableKmsVault(ctx),
			"oci_logging_log":                            tableLoggingLog(ctx),
			"oci_logging_log_group":                      tableLoggingLogGroup(ctx),
			"oci_objectstorage_bucket":                   tableObjectStorageBucket(ctx),
			"oci_ons_notification_topic":                 tableOnsNotificationTopic(ctx),
			"oci_ons_subscription":                       tableOnsSubscription(ctx),
			"oci_region":                                 tableIdentityRegion(ctx),
		},
	}
	return p
}
