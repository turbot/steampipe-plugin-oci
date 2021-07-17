## v0.0.15 [2021-07-16]

_What's new?_

- New tables added
  - [oci_budget_alert_rule](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_budget_alert_rule) ([#198](https://github.com/turbot/steampipe-plugin-oci/pull/198))
  - [oci_budget_budget](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_budget_budget) ([#197](https://github.com/turbot/steampipe-plugin-oci/pull/197))
  - [oci_core_public_ip](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_public_ip) ([#194](https://github.com/turbot/steampipe-plugin-oci/pull/194))
  - [oci_database_autonomous_database](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database) ([#162](https://github.com/turbot/steampipe-plugin-oci/pull/162))
  - [oci_database_mysql_db_system](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_mysql_db_system) ([#169](https://github.com/turbot/steampipe-plugin-oci/pull/169))
  - [oci_dns_rrset](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_dns_rrset) ([#112](https://github.com/turbot/steampipe-plugin-oci/pull/112))
  - [oci_mysql_channel](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_channel) ([#171](https://github.com/turbot/steampipe-plugin-oci/pull/171))
  - [oci_nosql_table](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table) ([#193](https://github.com/turbot/steampipe-plugin-oci/pull/193))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v131--2021-07-15) ([#196](https://github.com/turbot/steampipe-plugin-oci/pull/196))
- Updated: oci-go-sdk to v44 ([#206](https://github.com/turbot/steampipe-plugin-oci/pull/206))
- Updated: Add column `object_lifecycle_policy` to `oci_objectstorage_bucket` table ([#187](https://github.com/turbot/steampipe-plugin-oci/pull/187))
- Updated: Minor cleanup in docs/index.md ([#173](https://github.com/turbot/steampipe-plugin-oci/pull/173))

## v0.0.14 [2021-07-01]

_Enhancements_

- Updated: Improved example queries in `oci_kms_key_version` table doc ([#155](https://github.com/turbot/steampipe-plugin-oci/pull/155))

_Bug fixes_

- Fixed: Compartments in creating and deleting states should not be retrieved in `oci_multi_region` table ([#151](https://github.com/turbot/steampipe-plugin-oci/pull/151))
- Fixed: Remove unused `region` column in `oci_dns_zone` table ([#157](https://github.com/turbot/steampipe-plugin-oci/pull/157))
- Fixed: Remove unused `region` column in `oci_dns_tsig_key` table ([#159](https://github.com/turbot/steampipe-plugin-oci/pull/159))

## v0.0.13 [2021-06-24]

_What's new?_

- New tables added
  - [oci_identity_api_key](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_api_key) ([#143](https://github.com/turbot/steampipe-plugin-oci/pull/143))
  - [oci_kms_key](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_kms_key) ([#145](https://github.com/turbot/steampipe-plugin-oci/pull/145))
  - [oci_kms_key_version](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_kms_key_version) ([#149](https://github.com/turbot/steampipe-plugin-oci/pull/149))

_Bug fixes_

- Fixed: Example query in `oci_identity_compartment` docs ([#134](https://github.com/turbot/steampipe-plugin-oci/pull/134))

## v0.0.12 [2021-06-17]

_What's new?_

- New tables added
  - [oci_file_storage_snapshot](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_file_storage_snapshot) ([#121](https://github.com/turbot/steampipe-plugin-oci/pull/121))
  - [oci_identity_tag_default](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_tag_default) ([#141](https://github.com/turbot/steampipe-plugin-oci/pull/141))

## v0.0.11 [2021-06-10]

_What's new?_

- New tables added
  - [oci_apigateway_api](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_apigateway_api) ([#130](https://github.com/turbot/steampipe-plugin-oci/pull/130))
  - [oci_core_boot_volume_backup](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_backup) ([#129](https://github.com/turbot/steampipe-plugin-oci/pull/129))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.10](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v0210-2021-06-09)

## v0.0.10 [2021-06-03]

_What's new?_

- New tables added
  - [oci_cloud_guard_target](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_guard_target) ([#122](https://github.com/turbot/steampipe-plugin-oci/pull/122))
  - [oci_dns_tsig_key](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_dns_tsig_key) ([#119](https://github.com/turbot/steampipe-plugin-oci/pull/119))

## v0.0.9 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [oci_cloud_guard_managed_list](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_guard_managed_list) ([#78](https://github.com/turbot/steampipe-plugin-oci/pull/78))
  - [oci_dns_zone](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_dns_zone) ([#101](https://github.com/turbot/steampipe-plugin-oci/pull/101))
  - [oci_file_storage_file_system](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_file_storage_file_system) ([#116](https://github.com/turbot/steampipe-plugin-oci/pull/116))
  - [oci_identity_tenancy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_tenancy) ([#111](https://github.com/turbot/steampipe-plugin-oci/pull/111))

_Bug fixes_

- Fixed: Columns now hydrate properly in `oci_cloud_guard_detector_recipe` and `oci_cloud_guard_responder_recipe` tables ([#124](https://github.com/turbot/steampipe-plugin-oci/pull/124))

## v0.0.8 [2021-05-20]

_What's new?_

- New tables added
  - [oci_core_network_security_group](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_network_security_group) ([#95](https://github.com/turbot/steampipe-plugin-oci/pull/95))
  - [oci_identity_tag_namespace](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_tag_namespace) ([#110](https://github.com/turbot/steampipe-plugin-oci/pull/110))
  - [oci_logging_log](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_logging_log) ([#100](https://github.com/turbot/steampipe-plugin-oci/pull/100))

## v0.0.7 [2021-05-13]

_What's new?_

- New tables added
  - [oci_core_drg](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_drg) ([#90](https://github.com/turbot/steampipe-plugin-oci/pull/90))
  - [oci_core_volume_backup_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume_backup_policy) ([#99](https://github.com/turbot/steampipe-plugin-oci/pull/99))

_Enhancements_

- Updated: README.md and docs/index.md now contain links to our Slack community ([#114](https://github.com/turbot/steampipe-plugin-oci/pull/114))

_Bug fixes_

- Fixed: Remove connection name in `oci_events_rule` table doc example query ([#109](https://github.com/turbot/steampipe-plugin-oci/pull/109))

## v0.0.6 [2021-04-29]

_What's new?_

- New tables added
  - [oci_autoscaling_auto_scaling_configuration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_autoscaling_auto_scaling_configuration) ([#91](https://github.com/turbot/steampipe-plugin-oci/pull/91))
  - [oci_cloud_guard_configuration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_guard_configuration) ([#88](https://github.com/turbot/steampipe-plugin-oci/pull/88))
  - [oci_cloud_guard_detector_recipe](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_guard_detector_recipe) ([#74](https://github.com/turbot/steampipe-plugin-oci/pull/74))
  - [oci_cloud_guard_responder_recipe](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_guard_responder_recipe) ([#81](https://github.com/turbot/steampipe-plugin-oci/pull/81))
  - [oci_functions_application](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_functions_application) ([#82](https://github.com/turbot/steampipe-plugin-oci/pull/82))
  - [oci_kms_vault](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_kms_vault) ([#86](https://github.com/turbot/steampipe-plugin-oci/pull/86))
  - [oci_logging_log_group](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_logging_log_group) ([#76](https://github.com/turbot/steampipe-plugin-oci/pull/76))

## v0.0.5 [2021-04-22]

_What's new?_

- New tables added
  - [oci_core_dhcp_options](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_dhcp_options) ([#75](https://github.com/turbot/steampipe-plugin-oci/pull/75))
  - [oci_core_image](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_image) ([#15](https://github.com/turbot/steampipe-plugin-oci/pull/15))
  - [oci_core_nat_gateway](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_nat_gateway) ([#45](https://github.com/turbot/steampipe-plugin-oci/pull/45))
  - [oci_core_security_list](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_security_list) ([#65](https://github.com/turbot/steampipe-plugin-oci/pull/65))
  - [oci_core_service_gateway](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_service_gateway) ([#49](https://github.com/turbot/steampipe-plugin-oci/pull/49))
  - [oci_core_subnet](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_subnet) ([#52](https://github.com/turbot/steampipe-plugin-oci/pull/52))
  - [oci_events_rule](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_events_rule) ([#79](https://github.com/turbot/steampipe-plugin-oci/pull/79))
  - [oci_identity_auth_token](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_auth_token) ([#68](https://github.com/turbot/steampipe-plugin-oci/pull/68))
  - [oci_identity_customer_secret_key](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_customer_secret_key) ([#66](https://github.com/turbot/steampipe-plugin-oci/pull/66))
  - [oci_identity_dynamic_group](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_dynamic_group) ([#34](https://github.com/turbot/steampipe-plugin-oci/pull/34))
  - [oci_ons_notification_topic](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ons_notification_topic) ([#37](https://github.com/turbot/steampipe-plugin-oci/pull/37))
  - [oci_ons_subscription](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ons_subscription) ([#43](https://github.com/turbot/steampipe-plugin-oci/pull/43))

## v0.0.4 [2021-04-15]

_What's new?_

- New tables added
  - [oci_core_local_peering_gateway](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_local_peering_gateway) ([#58](https://github.com/turbot/steampipe-plugin-oci/pull/58))
  - [oci_core_vcn](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_vcn) ([#23](https://github.com/turbot/steampipe-plugin-oci/pull/23))
  - [oci_core_volume_backup](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume_backup) ([#31](https://github.com/turbot/steampipe-plugin-oci/pull/31))
  - [oci_identity_network_source](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_network_source) ([#39](https://github.com/turbot/steampipe-plugin-oci/pull/39))

## v0.0.3 [2021-04-08]

_What's new?_

- New tables added
  - [oci_core_internet_gateway](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_internet_gateway) ([#21](https://github.com/turbot/steampipe-plugin-oci/pull/21))
  - [oci_core_route_table](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_route_table) ([#44](https://github.com/turbot/steampipe-plugin-oci/pull/44))
  - [oci_core_volume](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume) ([#10](https://github.com/turbot/steampipe-plugin-oci/pull/10))

## v0.0.2 [2021-04-02]

_Bug fixes_

- Fixed: `Table definitions & examples` link now points to the correct location ([#35](https://github.com/turbot/steampipe-plugin-oci/pull/35))

## v0.0.1 [2021-04-01]

_What's new?_

- New tables added
  - [oci_core_instance](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_instance) ([#1](https://github.com/turbot/steampipe-plugin-oci/pull/1))
  - [oci_identity_authentication_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_authentication_policy) ([#14](https://github.com/turbot/steampipe-plugin-oci/pull/14))
  - [oci_identity_compartment](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_compartment) ([#14](https://github.com/turbot/steampipe-plugin-oci/pull/14))
  - [oci_identity_group](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_group) ([#8](https://github.com/turbot/steampipe-plugin-oci/pull/8))
  - [oci_identity_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_policy) ([#9](https://github.com/turbot/steampipe-plugin-oci/pull/9))
  - [oci_identity_user](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_user) ([#1](https://github.com/turbot/steampipe-plugin-oci/pull/1))
  - [oci_objectstorage_bucket](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_objectstorage_bucket) ([#17](https://github.com/turbot/steampipe-plugin-oci/pull/17))
  - [oci_region](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_region) ([#1](https://github.com/turbot/steampipe-plugin-oci/pull/1))
