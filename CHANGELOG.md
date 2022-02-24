## v0.8.1 [2022-02-24]

_Bug fixes_

- Fixed the `oci_kms_key` table to list keys from different compartments instead of fetching keys from only the vault's compartment ([#365](https://github.com/turbot/steampipe-plugin-oci/pull/365))
- The following tables have been renamed to shorten their names to avoid hitting a PostgreSQL identifier character limit that would prevent proper schema generation. Any scripts or workflows that use should be updated to use the updated names. ([#364](https://github.com/turbot/steampipe-plugin-oci/pull/364))
  - `oci_database_autonomous_database_metric_cpu_utilization` renamed to `oci_database_autonomous_db_metric_cpu_utilization`
  - `oci_database_autonomous_database_metric_cpu_utilization_daily` renamed to `oci_database_autonomous_db_metric_cpu_utilization_daily`
  - `oci_database_autonomous_database_metric_cpu_utilization_hourly` renamed to `oci_database_autonomous_db_metric_cpu_utilization_hourly`
  - `oci_database_autonomous_database_metric_storage_utilization` renamed to `oci_database_autonomous_db_metric_storage_utilization`
  - `oci_database_autonomous_database_metric_storage_utilization_daily` renamed to `oci_database_autonomous_db_metric_storage_utilization_daily`
  - `oci_database_autonomous_database_metric_storage_utilization_hourly` renamed to `oci_database_autonomous_db_metric_storage_utilization_hourly`

## v0.8.0 [2022-02-18]

_What's new?_

- New tables added
  - [oci_core_vnic_attachment](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_vnic_attachment) ([#346](https://github.com/turbot/steampipe-plugin-oci/pull/346))
  - [oci_database_autonomous_database_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_cpu_utilization_hourly) ([#353](https://github.com/turbot/steampipe-plugin-oci/pull/353))
  - [oci_database_autonomous_database_metric_storage_utilization_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_storage_utilization_hourly) ([#359](https://github.com/turbot/steampipe-plugin-oci/pull/359))
  - [oci_mysql_db_system_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_cpu_utilization_hourly) ([#352](https://github.com/turbot/steampipe-plugin-oci/pull/352))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.0.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v203--2022-02-14) ([#354](https://github.com/turbot/steampipe-plugin-oci/pull/354))

_Bug fixes_

- Fixed the lower limit in the `oci_dns_rrset` table ([#357](https://github.com/turbot/steampipe-plugin-oci/pull/357))
- Updated the column type  of `time_created` column to `TIMESTAMP` in `oci_objectstorage_bucket` table ([#348](https://github.com/turbot/steampipe-plugin-oci/pull/348))

## v0.7.0 [2022-01-19]

_What's new?_

- New tables added
  - [oci_objectstorage_object](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_objectstorage_object) ([#342](https://github.com/turbot/steampipe-plugin-oci/pull/342))
  - [oci_vault_secret](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_vault_secret) ([#343](https://github.com/turbot/steampipe-plugin-oci/pull/343))

## v0.6.0 [2022-01-12]

_What's new?_

- New tables added
  - [oci_functions_function](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_functions_function) ([#340](https://github.com/turbot/steampipe-plugin-oci/pull/340))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#339](https://github.com/turbot/steampipe-plugin-oci/pull/339))

_Bug fixes_

- Fixed the Cloud Guard tables to correctly return the result instead of an empty table ([#335](https://github.com/turbot/steampipe-plugin-oci/pull/335))

## v0.5.0 [2021-11-24]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#331](https://github.com/turbot/steampipe-plugin-oci/pull/331))
- Recompiled plugin with Go version 1.17 ([#332](https://github.com/turbot/steampipe-plugin-oci/pull/332))

_Bug fixes_

- `oci_apigateway_api` table will now return data available in each region configured in the plugin config (oci.spc) file ([#328](https://github.com/turbot/steampipe-plugin-oci/pull/328))
- `oci_logging_log_group` table will no longer return `InvalidParameter` error in get call ([#327](https://github.com/turbot/steampipe-plugin-oci/pull/327))
- `oci_core_nat_gateway` table will no longer return `400` or `404` error in get call ([#327](https://github.com/turbot/steampipe-plugin-oci/pull/327))

## v0.4.0 [2021-11-03]

_Enhancements_

- Updated: Add additional optional key quals, filter support, and context cancellation handling and improve hydrate with cache functionality across all the  tables ([#306](https://github.com/turbot/steampipe-plugin-oci/pull/306)) ([#317](https://github.com/turbot/steampipe-plugin-oci/pull/317))
- Updated: Add `title` column in `oci_ons_subscription` table ([#313](https://github.com/turbot/steampipe-plugin-oci/pull/313))

## v0.3.1 [2021-09-13]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v151--2021-09-13) ([#309](https://github.com/turbot/steampipe-plugin-oci/pull/309))

## v0.3.0 [2021-09-02]

_What's new?_

- Added support for the region `sa-vinhedo-1` across all the tables ([#302](https://github.com/turbot/steampipe-plugin-oci/pull/302))

## v0.2.0 [2021-08-26]

_What's new?_

- New tables added
  - [oci_analytics_instance](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_analytics_instance) ([#280](https://github.com/turbot/steampipe-plugin-oci/pull/280))
  - [oci_core_image_custom](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_image_custom) ([#276](https://github.com/turbot/steampipe-plugin-oci/pull/276))
  - [oci_database_db](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_db) ([#287](https://github.com/turbot/steampipe-plugin-oci/pull/287))
  - [oci_database_db_home](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_db_home) ([#288](https://github.com/turbot/steampipe-plugin-oci/pull/288))
  - [oci_database_db_system](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_db_system) ([#286](https://github.com/turbot/steampipe-plugin-oci/pull/286))
  - [oci_database_software_image](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_software_image) ([#291](https://github.com/turbot/steampipe-plugin-oci/pull/291))
  - [oci_mysql_configuration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_configuration) ([#255](https://github.com/turbot/steampipe-plugin-oci/pull/255))
  - [oci_mysql_configuration_custom](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_configuration_custom) ([#285](https://github.com/turbot/steampipe-plugin-oci/pull/285))
  - [oci_resource_search](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_resource_search) ([#148](https://github.com/turbot/steampipe-plugin-oci/pull/148))

_Enhancements_

- Updated: `oci_core_instance` table now includes `shape_config_max_vnic_attachments`, `shape_config_memory_in_gbs`, `shape_config_networking_bandwidth_in_gbps`, `shape_config_ocpus`, `shape_config_baseline_ocpu_utilization`, `shape_config_gpus`, `shape_config_local_disks`, `shape_config_local_disks_total_size_in_gbs` columns ([#294](https://github.com/turbot/steampipe-plugin-oci/pull/294))

_Bug fixes_

- Fixed: `oci_core_image` table no longer includes duplicate data and now only lists platform (standard) images ([#275](https://github.com/turbot/steampipe-plugin-oci/pull/275))

## v0.1.0 [2021-08-06]

_What's new?_

- New tables added
  - [oci_core_public_ip_pool](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_public_ip_pool) ([#260](https://github.com/turbot/steampipe-plugin-oci/pull/260))
  - [oci_file_storage_mount_target](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_file_storage_mount_target) ([#257](https://github.com/turbot/steampipe-plugin-oci/pull/257))

_Bug fixes_

- Fixed: Restrict get API calls to the root compartment and one zone per region in several core and file storage tables ([#264](https://github.com/turbot/steampipe-plugin-oci/pull/264))

## v0.0.17 [2021-07-31]

_What's new?_

- New tables added
  - [oci_core_boot_volume_replica](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_replica) ([#236](https://github.com/turbot/steampipe-plugin-oci/pull/236))
  - [oci_core_load_balancer](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_load_balancer) ([#241](https://github.com/turbot/steampipe-plugin-oci/pull/241))
  - [oci_database_autonomous_database_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_cpu_utilization) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_database_autonomous_database_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_cpu_utilization_daily) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_database_autonomous_database_metric_storage_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_storage_utilization) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_database_autonomous_database_metric_storage_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_autonomous_database_metric_storage_utilization_daily) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_connections](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_connections) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_connections_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_connections_daily) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_cpu_utilization) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_cpu_utilization_daily) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_memory_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_memory_utilization) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))
  - [oci_mysql_db_system_metric_memory_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_memory_utilization_daily) ([#254](https://github.com/turbot/steampipe-plugin-oci/pull/254))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v141--2021-07-20) ([#256](https://github.com/turbot/steampipe-plugin-oci/pull/256))
- Updated: Add integration test for `oci_core_instance` table ([#261](https://github.com/turbot/steampipe-plugin-oci/pull/261))

_Bug fixes_

- Fixed: Cache keys for monitoring service and identity service regional connection information are now correct ([#259](https://github.com/turbot/steampipe-plugin-oci/pull/259))
- Fixed: Rename `table_core_volume_backup.go` to `table_oci_core_volume_backup.go` ([#245](https://github.com/turbot/steampipe-plugin-oci/pull/245))

## v0.0.16 [2021-07-22]

_What's new?_

- New tables added
  - [oci_core_block_volume_replica](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_block_volume_replica) ([#202](https://github.com/turbot/steampipe-plugin-oci/pull/202))
  - [oci_core_boot_volume](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume) ([#208](https://github.com/turbot/steampipe-plugin-oci/pull/208))
  - [oci_core_boot_volume_attachment](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_attachment) ([#223](https://github.com/turbot/steampipe-plugin-oci/pull/223))
  - [oci_core_boot_volume_metric_read_ops](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_read_ops) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_boot_volume_metric_read_ops_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_read_ops_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_boot_volume_metric_read_ops_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_read_ops_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_boot_volume_metric_write_ops](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_write_ops) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_boot_volume_metric_write_ops_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_write_ops_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_boot_volume_metric_write_ops_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_boot_volume_metric_write_ops_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_instance_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_instance_metric_cpu_utilization) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_instance_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_instance_metric_cpu_utilization_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_instance_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_instance_metric_cpu_utilization_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_core_network_load_balancer](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_network_load_balancer) ([#224](https://github.com/turbot/steampipe-plugin-oci/pull/224))
  - [oci_core_volume_attachment](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume_attachment) ([#207](https://github.com/turbot/steampipe-plugin-oci/pull/207))
  - [oci_identity_availability_domain](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_availability_domain) ([#210](https://github.com/turbot/steampipe-plugin-oci/pull/210))
  - [oci_mysql_backup](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_backup) ([#170](https://github.com/turbot/steampipe-plugin-oci/pull/170))
  - [oci_nosql_table_metric_read_throttle_count](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_read_throttle_count) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_read_throttle_count_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_read_throttle_count_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_read_throttle_count_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_read_throttle_count_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_storage_utilization](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_storage_utilization) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_storage_utilization_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_storage_utilization_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_storage_utilization_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_storage_utilization_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_write_throttle_count](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_write_throttle_count) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_write_throttle_count_daily](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_write_throttle_count_daily) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))
  - [oci_nosql_table_metric_write_throttle_count_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_nosql_table_metric_write_throttle_count_hourly) ([#204](https://github.com/turbot/steampipe-plugin-oci/pull/204))

_Enhancements_

- Updated: Add column `region` to all metric tables ([#240](https://github.com/turbot/steampipe-plugin-oci/pull/240))
- Updated: Add column `region` to `oci_mysql_db_system` table ([#235](https://github.com/turbot/steampipe-plugin-oci/pull/235))

_Bug fixes_

- Fixed: Network load balancer service connection no longer fails due to undeclared tenant ID ([#232](https://github.com/turbot/steampipe-plugin-oci/pull/232))

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
