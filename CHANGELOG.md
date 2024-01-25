## v0.35.0 [2024-01-25]

_What's new?_

- New tables added
  - [oci_identity_db_credential](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_db_credential) ([#596](https://github.com/turbot/steampipe-plugin-oci/pull/596))

## v0.34.0 [2024-01-22]

_What's new?_

- New tables added
  - [oci_identity_domain](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_identity_domain) ([#587](https://github.com/turbot/steampipe-plugin-oci/pull/587))
  - [oci_database_cloud_vm_cluster](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_cloud_vm_cluster) ([#583](https://github.com/turbot/steampipe-plugin-oci/pull/583))

_Enhancements_

- Added the `snapshot_time` and `snapshot_type` columns to `oci_file_storage_snapshot` table. ([#592](https://github.com/turbot/steampipe-plugin-oci/pull/592))
- Added the `kms_key_version_id`, `vault_id`, `sid_prefix` and `is_cdb` columns to `oci_database_db` table. ([#591](https://github.com/turbot/steampipe-plugin-oci/pull/591))

## v0.33.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#575](https://github.com/turbot/steampipe-plugin-oci/pull/575))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#575](https://github.com/turbot/steampipe-plugin-oci/pull/575))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-oci/blob/main/docs/LICENSE). ([#575](https://github.com/turbot/steampipe-plugin-oci/pull/575))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#574](https://github.com/turbot/steampipe-plugin-oci/pull/574))

## v0.32.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#553](https://github.com/turbot/steampipe-plugin-oci/pull/553))

## v0.32.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#549](https://github.com/turbot/steampipe-plugin-oci/pull/549))
- Recompiled plugin with Go version `1.21`. ([#549](https://github.com/turbot/steampipe-plugin-oci/pull/549))

## v0.31.0 [2023-09-27]

_Enhancements_

- Added the `last_successful_login_time` column to `oci_identity_user` table. ([#547](https://github.com/turbot/steampipe-plugin-oci/pull/547))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v551-2023-07-26). ([#542](https://github.com/turbot/steampipe-plugin-oci/pull/542))

## v0.30.1 [2023-09-01]

_Bug fixes_

- Fixed the `index out of range` error while querying `oci_database_autonomous_database`, `oci_events_rule` and `oci_identity_tag_default` tables, when the `regions` config arg is not set in the `oci.spc` file. ([#537](https://github.com/turbot/steampipe-plugin-oci/pull/537))

## v0.30.0 [2023-08-22]

_Enhancements_

- Added example queries in the `oci_nosql_table` table doc for counting the number of child tables of a parent table. ([#535](https://github.com/turbot/steampipe-plugin-oci/pull/535))

_Bug fixes_

- Fixed the validation of the `regions` argument in the configuration file to generate an error in case an unsupported or an invalid region is provided. ([#534](https://github.com/turbot/steampipe-plugin-oci/pull/534))

## v0.29.0 [2023-07-14]

_What's new?_

- New tables added
  - [oci_container_instances_container_instance](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_container_instances_container_instance) ([#525](https://github.com/turbot/steampipe-plugin-oci/pull/525))
  - [oci_container_instances_container](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_container_instances_container) ([#526](https://github.com/turbot/steampipe-plugin-oci/pull/526))
  - [oci_logging_search](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_logging_search) ([#527](https://github.com/turbot/steampipe-plugin-oci/pull/527))

_Enhancements_

- Updated the `docs/index.md`` file to include multi-tenant configuration examples. ([#531](https://github.com/turbot/steampipe-plugin-oci/pull/531))

## v0.28.0 [2023-06-29]

_What's new?_

- New tables added
  - [oci_devops_project](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_devops_project) ([#523](https://github.com/turbot/steampipe-plugin-oci/pull/523)) (Thanks [lucasjellema](https://github.com/lucasjellema) for the contribution!)
  - [oci_devops_repository](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_devops_repository) ([#522](https://github.com/turbot/steampipe-plugin-oci/pull/522))

## v0.27.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#520](https://github.com/turbot/steampipe-plugin-oci/pull/520))

## v0.26.0 [2023-06-13]

_What's new?_

- New tables added: (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)
  - [oci_adm_knowledge_base](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_adm_knowledge_base) ([#509](https://github.com/turbot/steampipe-plugin-oci/pull/509))
  - [oci_adm_vulnerability_audit](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_adm_vulnerability_audit) ([#509](https://github.com/turbot/steampipe-plugin-oci/pull/509))
  - [oci_application_migration_migration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_application_migration_migration) ([#510](https://github.com/turbot/steampipe-plugin-oci/pull/510))
  - [oci_application_migration_source](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_application_migration_source) ([#510](https://github.com/turbot/steampipe-plugin-oci/pull/510))
  - [oci_autoscaling_auto_scaling_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_autoscaling_auto_scaling_policy) ([#513](https://github.com/turbot/steampipe-plugin-oci/pull/513))
  - [oci_bds_bds_instance](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_bds_bds_instance) ([#514](https://github.com/turbot/steampipe-plugin-oci/pull/514))
  - [oci_certificate_authority_bundle](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_authority_bundle) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_association](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_association) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_ca_bundle](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_ca_bundle) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_certificate](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_certificate) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_certificate_authority](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_certificate_authority) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_certificate_authority_version](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_certificate_authority_version) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))
  - [oci_certificate_management_certificate_version](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_certificate_management_certificate_version) ([#516](https://github.com/turbot/steampipe-plugin-oci/pull/516))

_Bug fixes_

- Fixed `oci_identity_policy` and `oci_identity_availability_domain` tables to return data for all compartments instead of only `root` compartment. ([#505](https://github.com/turbot/steampipe-plugin-oci/pull/505)) ([#512](https://github.com/turbot/steampipe-plugin-oci/pull/512))

## v0.25.0 [2023-06-03]

_What's new?_

- New tables added: (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)
  - [oci_ai_anomaly_detection_ai_private_endpoint](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ai_anomaly_detection_ai_private_endpoint) ([#508](https://github.com/turbot/steampipe-plugin-oci/pull/508))
  - [oci_ai_anomaly_detection_data_asset](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ai_anomaly_detection_data_asset) ([#508](https://github.com/turbot/steampipe-plugin-oci/pull/508))
  - [oci_ai_anomaly_detection_model](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ai_anomaly_detection_model) ([#508](https://github.com/turbot/steampipe-plugin-oci/pull/508))
  - [oci_ai_anomaly_detection_project](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_ai_anomaly_detection_project) ([#508](https://github.com/turbot/steampipe-plugin-oci/pull/508))
  - [oci_artifacts_container_image_signature](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_artifacts_container_image_signature) ([#511](https://github.com/turbot/steampipe-plugin-oci/pull/511))
  - [oci_artifacts_container_image](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_artifacts_container_image) ([#511](https://github.com/turbot/steampipe-plugin-oci/pull/511))
  - [oci_artifacts_container_repository](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_artifacts_container_repository) ([#511](https://github.com/turbot/steampipe-plugin-oci/pull/511))
  - [oci_artifacts_generic_artifact](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_artifacts_generic_artifact) ([#511](https://github.com/turbot/steampipe-plugin-oci/pull/511))
  - [oci_artifacts_repository](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_artifacts_repository) ([#511](https://github.com/turbot/steampipe-plugin-oci/pull/511))
  - [oci_network_firewall_firewall](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_network_firewall_firewall) ([#507](https://github.com/turbot/steampipe-plugin-oci/pull/507))
  - [oci_network_firewall_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_network_firewall_policy) ([#507](https://github.com/turbot/steampipe-plugin-oci/pull/507))

## v0.24.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#503](https://github.com/turbot/steampipe-plugin-oci/pull/503))

## v0.23.0 [2023-04-14]

_Enhancements_

- Added column `tenant_name` across all the tables. ([#495](https://github.com/turbot/steampipe-plugin-oci/pull/495))

## v0.22.0 [2023-04-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which adds go-getter support to dynamic tables. ([#496](https://github.com/turbot/steampipe-plugin-oci/pull/496))

## v0.21.0 [2023-03-15]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#492](https://github.com/turbot/steampipe-plugin-oci/pull/492))

## v0.20.0 [2023-03-03]

_What's new?_

- New tables added
  - [oci_bastion_bastion](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_bastion_bastion) ([#476](https://github.com/turbot/steampipe-plugin-oci/pull/476)) (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)
  - [oci_bastion_session](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_bastion_session) ([#476](https://github.com/turbot/steampipe-plugin-oci/pull/476)) (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)
  - [oci_queue_queue](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_queue_queue) ([#481](https://github.com/turbot/steampipe-plugin-oci/pull/481)) (Thanks [@lucasjellema](https://github.com/lucasjellema) for the contribution!)

_Enhancements_

- Added columns `cpu_core_count` and `memory_size_in_gbs` to `oci_mysql_db_system` table. ([#483](https://github.com/turbot/steampipe-plugin-oci/pull/483)) (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)

## v0.19.1 [2023-02-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#477](https://github.com/turbot/steampipe-plugin-oci/pull/477))

## v0.19.0 [2023-01-30]

_What's new?_

- New tables added
  - [oci_core_cluster_network](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_cluster_network) ([#452](https://github.com/turbot/steampipe-plugin-oci/pull/452))
  - [oci_core_instance_configuration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_instance_configuration) ([#464](https://github.com/turbot/steampipe-plugin-oci/pull/464))
  - [oci_core_volume_default_backup_policy](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume_default_backup_policy) ([#463](https://github.com/turbot/steampipe-plugin-oci/pull/463))
  - [oci_core_volume_group](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_core_volume_group) ([#467](https://github.com/turbot/steampipe-plugin-oci/pull/467))

_Enhancements_

- Added columns `kms_key_id`, `kms_key_version_id` and `vault_id` to `oci_database_autonomous_database` table. ([#469](https://github.com/turbot/steampipe-plugin-oci/pull/469))

_Bug fixes_

- Fixed the `exports` column in `oci_file_storage_file_system` to correctly return data instead of `null` when an `id` is passed in the `where` clause. ([#466](https://github.com/turbot/steampipe-plugin-oci/pull/466))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.11](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4111-2023-01-24) which fixes the issue of non-caching of all the columns of the queried table. ([#474](https://github.com/turbot/steampipe-plugin-oci/pull/474))

## v0.18.0 [2023-01-06]

_Enhancements_

- Added column `exports` to the `oci_file_storage_file_system` table. ([#450](https://github.com/turbot/steampipe-plugin-oci/pull/450))
- Added an example query in the `oci_identity_compartment` table that lists out the full path of OCI compartments. ([#428](https://github.com/turbot/steampipe-plugin-oci/pull/428)) (Thanks [@AnykeyNL](https://github.com/AnykeyNL) for the contribution!)

_Dependencies_

- Recompiled plugin with [oci-go-sdk v65.28.0](https://github.com/oracle/oci-go-sdk/blob/master/CHANGELOG.md#65280---2022-12-13). ([#433](https://github.com/turbot/steampipe-plugin-oci/pull/433)) (Thanks [@scotti-fletcher](https://github.com/scotti-fletcher) for the contribution!)

## v0.17.2 [2022-11-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#426](https://github.com/turbot/steampipe-plugin-oci/pull/426))

## v0.17.1 [2022-09-19]

_Bug fixes_

- Fixed `delivery_policy` column not returning data in `oci_ons_subscription` table when specifying the `id` column. ([#370](https://github.com/turbot/steampipe-plugin-oci/pull/370))

## v0.17.0 [2022-09-06]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v416-2022-09-02) which includes several caching and memory management improvements. ([#422](https://github.com/turbot/steampipe-plugin-oci/pull/422))
- Recompiled plugin with Go version `1.19`. ([#422](https://github.com/turbot/steampipe-plugin-oci/pull/422))

## v0.16.0 [2022-07-14]

_Breaking changes_

- Fixed the typo in the table name to use `oci_database_autonomous_db_metric_storage_utilization` instead of `oci_database_autonomous_dd_metric_storage_utilization`. ([#417](https://github.com/turbot/steampipe-plugin-oci/pull/417))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#418](https://github.com/turbot/steampipe-plugin-oci/pull/418))

## v0.15.0 [2022-07-01]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v331--2022-06-30). ([#412](https://github.com/turbot/steampipe-plugin-oci/pull/412))

## v0.14.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#407](https://github.com/turbot/steampipe-plugin-oci/pull/407))

## v0.14.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#401](https://github.com/turbot/steampipe-plugin-oci/pull/401))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#400](https://github.com/turbot/steampipe-plugin-oci/pull/400))

## v0.13.0 [2022-04-13]

_What's new?_

- Added optional config arguments `max_error_retry_attempts` and `min_error_retry_delay` to allow customization of the error retry timings. For more information please see [OCI plugin configuration](https://hub.steampipe.io/plugins/turbot/oci#configuration). ([#397](https://github.com/turbot/steampipe-plugin-oci/pull/397))

## v0.12.0 [2022-04-06]

_What's new?_

- New tables added
  - [oci_containerengine_cluster](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_containerengine_cluster) ([#386](https://github.com/turbot/steampipe-plugin-oci/pull/386))
  - [oci_mysql_heat_wave_cluster](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_heat_wave_cluster) ([#385](https://github.com/turbot/steampipe-plugin-oci/pull/385))
  - [oci_resourcemanager_stack](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_resourcemanager_stack) ([#387](https://github.com/turbot/steampipe-plugin-oci/pull/387))

## v0.11.0 [2022-03-30]

_What's new?_

- New tables added
  - [oci_database_pluggable_database](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_database_pluggable_database) ([#384](https://github.com/turbot/steampipe-plugin-oci/pull/384))
  - [oci_mysql_db_system_metric_connections_hourly](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_mysql_db_system_metric_connections_hourly) ([#382](https://github.com/turbot/steampipe-plugin-oci/pull/382))
  - [oci_streaming_stream](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_streaming_stream) ([#381](https://github.com/turbot/steampipe-plugin-oci/pull/381))

_Enhancements_

- Added column `ipv6_cidr_blocks` to the `oci_core_vcn` table ([#390](https://github.com/turbot/steampipe-plugin-oci/pull/390))

## v0.10.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#377](https://github.com/turbot/steampipe-plugin-oci/pull/377))

## v0.9.0 [2022-03-10]

_Enhancements_

- Added `volume_backup_policy_id` and `volume_backup_policy_assignment_id` columns to `oci_core_boot_volume` table ([#371](https://github.com/turbot/steampipe-plugin-oci/pull/371))
- Added `volume_backup_policy_id` and `volume_backup_policy_assignment_id` columns to `oci_core_volume` table ([#371](https://github.com/turbot/steampipe-plugin-oci/pull/371))

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
- Updated the column type of `time_created` column to `TIMESTAMP` in `oci_objectstorage_bucket` table ([#348](https://github.com/turbot/steampipe-plugin-oci/pull/348))

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

- Updated: Add additional optional key quals, filter support, and context cancellation handling and improve hydrate with cache functionality across all the tables ([#306](https://github.com/turbot/steampipe-plugin-oci/pull/306)) ([#317](https://github.com/turbot/steampipe-plugin-oci/pull/317))
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
