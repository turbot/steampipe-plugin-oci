---
title: "Steampipe Table: oci_application_migration_migration - Query OCI Application Migration Service Migrations using SQL"
description: "Allows users to query data related to migrations in Oracle Cloud Infrastructure's Application Migration Service."
---

# Table: oci_application_migration_migration - Query OCI Application Migration Service Migrations using SQL

**Deprecated. Use [oci_cloud_migrations_migration](https://hub.steampipe.io/plugins/turbot/oci/tables/oci_cloud_migrations_migration) instead.**

Oracle Cloud Infrastructure's Application Migration Service simplifies the migration of applications from on-premise data centers or other clouds to Oracle Cloud Infrastructure. It supports a wide range of source applications, including Java EE, Oracle WebLogic Server, Oracle SOA Suite, and more. The service provides a comprehensive solution for migrating applications, databases, and associated configurations.

## Table Usage Guide

The `oci_application_migration_migration` table provides insights into the migrations performed using Oracle Cloud Infrastructure's Application Migration Service. As a cloud engineer or database administrator, you can explore migration-specific details through this table, including migration status, type of source application, and associated metadata. Utilize it to monitor the progress of migrations, identify any issues, and ensure successful migration of applications to Oracle Cloud Infrastructure.

## Examples

### Basic info
Explore the status and details of your application migrations to understand their progress and configuration. This can help in managing and tracking the migration process effectively.

```sql+postgres
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration;
```

### List all inactive migrations
Discover the segments that consist of inactive migrations in your application to better manage resources and prioritize tasks. This can be particularly useful in maintaining efficiency and ensuring smooth operations within your system.

```sql+postgres
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  lifecycle_state <> 'ACTIVE';
```

### List all migrations created in the last 30 days
Identify recent application migrations within the past month. This can provide insights into the migration trends and help in assessing the migration workload.

```sql+postgres
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  time_created >= datetime('now', '-30 day');
```

### List source details of a particular migration
This query helps in assessing the details of a specific migration source in an application migration scenario. It's useful when you need to understand the lifecycle, account details, and regional information of the source associated with a particular migration.

```sql+postgres
select
  m.id as migration_id,
  m.display_name as migration_name,
  s.id as source_id,
  s.display_name as source_name,
  s.lifecycle_details,
  s.lifecycle_state as state,
  s.source_details ->> 'computeAccount' as source_account_name,
  s.source_details ->> 'region' as source_account_region,
  s.source_details ->> 'type' as source_account_type
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and m.display_name = 'migration-1';
```

```sql+sqlite
select
  m.id as migration_id,
  m.display_name as migration_name,
  s.id as source_id,
  s.display_name as source_name,
  s.lifecycle_details,
  s.lifecycle_state as state,
  json_extract(s.source_details, '$.computeAccount') as source_account_name,
  json_extract(s.source_details, '$.region') as source_account_region,
  json_extract(s.source_details, '$.type') as source_account_type
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and m.display_name = 'migration-1';
```

### List successfully migrated applications
Discover the segments that have successfully migrated applications in your OCI environment. This can help you track your migration progress and ensure that all applications have been transferred correctly.

```sql+postgres
select
  id as migration_id,
  display_name as migration_name,
  application_config,
  application_name,
  application_type
from
  oci_application_migration_migration
where
  migration_state = 'MIGRATION_SUCCEEDED';
```

```sql+sqlite
select
  id as migration_id,
  display_name as migration_name,
  application_config,
  application_name,
  application_type
from
  oci_application_migration_migration
where
  migration_state = 'MIGRATION_SUCCEEDED';
```

### List application migrations that have selective-migration enabled
Determine the areas in which applications have selective migration enabled. This query is useful for understanding which applications are set to migrate only specific parts, allowing for targeted and efficient migration processes.

```sql+postgres
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  is_selective_migration;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state
from
  oci_application_migration_migration
where
  is_selective_migration;
```