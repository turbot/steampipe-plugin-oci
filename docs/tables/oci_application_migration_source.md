---
title: "Steampipe Table: oci_application_migration_source - Query OCI Application Migration Sources using SQL"
description: "Allows users to query data from Oracle Cloud Infrastructure's Application Migration Sources."
---

# Table: oci_application_migration_source - Query OCI Application Migration Sources using SQL

**Deprecated. Use [azuread_group](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_group) instead.**

Oracle Cloud Infrastructure's Application Migration service helps you migrate applications from any environment to Oracle Cloud Infrastructure. It simplifies the migration process, reduces downtime, and ensures that your applications are migrated with minimal disruption. This service is beneficial for those who want to take advantage of Oracle Cloud Infrastructure's benefits without the hassle of manually migrating their applications.

## Table Usage Guide

The `oci_application_migration_source` table provides insights into the sources within Oracle Cloud Infrastructure's Application Migration service. As a Cloud engineer, explore source-specific details through this table, including the source type, lifecycle state, and associated metadata. Utilize it to uncover information about sources, such as those in different lifecycle states, the type of sources, and the verification of source details.

## Examples

### Basic info
Explore the basic information of your migration sources in Oracle Cloud Infrastructure (OCI) to understand their current lifecycle state and source details. This can be useful for tracking the status and origin of your migration sources, aiding in effective resource management and migration planning.

```sql+postgres
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  json_extract(source_details, '$.computeAccount') as source_account_name,
  json_extract(source_details, '$.region') as source_account_region,
  json_extract(source_details, '$.type') as source_account_type
from
  oci_application_migration_source;
```

### List all inactive sources
Discover the segments that are not currently active in your application migration process. This can be useful for identifying potential issues or bottlenecks, allowing you to optimize the migration process.

```sql+postgres
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  json_extract(source_details, '$.computeAccount') as source_account_name,
  json_extract(source_details, '$.region') as source_account_region,
  json_extract(source_details, '$.type') as source_account_type
from
  oci_application_migration_source
where
  lifecycle_state <> 'ACTIVE';
```

### List all sources created in the last 30 days
Explore which sources have been created in the last 30 days to better manage and understand the status of your application migration. This allows for efficient tracking of newly added sources, essential for maintaining an updated and secure application environment.

```sql+postgres
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  json_extract(source_details, '$.computeAccount') as source_account_name,
  json_extract(source_details, '$.region') as source_account_region,
  json_extract(source_details, '$.type') as source_account_type
from
  oci_application_migration_source
where
  time_created >= datetime('now', '-30 day');
```

### List all migrations under a particular source
Explore which migrations are associated with a specific source in order to manage and track application migration processes. This is particularly useful for identifying the state and details of migrations, aiding in the organization and monitoring of migration tasks.

```sql+postgres
select
  s.id as source_id,
  s.display_name as source_name,
  m.id as migration_id,
  m.display_name as migration_name,
  m.lifecycle_details,
  m.lifecycle_state as state
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and s.display_name = 'source-1';
```

```sql+sqlite
select
  s.id as source_id,
  s.display_name as source_name,
  m.id as migration_id,
  m.display_name as migration_name,
  m.lifecycle_details,
  m.lifecycle_state as state
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and s.display_name = 'source-1';
```