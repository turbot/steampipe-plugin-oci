---
title: "Steampipe Table: oci_cloud_migrations_migration - Query OCI Cloud Migrations Migrations using SQL"
description: "Allows users to query OCI Cloud Migrations Migrations."
---

# Table: oci_cloud_migrations_migration - Query OCI Cloud Migrations Migrations using SQL

Oracle Cloud Infrastructure (OCI) Cloud Migrations service helps enterprises to migrate their on-premises or other cloud workloads to OCI. It provides a centralized platform for planning, executing, and tracking migrations to OCI. The Cloud Migrations Migration resource represents a migration project that groups the assets, plans, and schedules for migrating workloads.

## Table Usage Guide

The `oci_cloud_migrations_migration` table provides insights into migrations within OCI Cloud Migrations service. As a cloud administrator or migration engineer, explore migration-specific details through this table, including state, creation time, and associated resources. Utilize it to track ongoing migrations, verify completion status, and monitor migration projects.

## Examples

### Basic info
Explore the basic details of your cloud migration projects to understand their current state and creation timeline. This can help with tracking migration progress and identifying any stalled or recently created projects.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  is_completed
from
  oci_cloud_migrations_migration;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  is_completed
from
  oci_cloud_migrations_migration;
```

### List active migrations
Identify migrations that are currently active in your environment to track ongoing migration projects and their creation dates. This helps in monitoring active workload transfers to OCI.

```sql+postgres
select
  id,
  display_name,
  time_created,
  is_completed
from
  oci_cloud_migrations_migration
where
  lifecycle_state = 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  time_created,
  is_completed
from
  oci_cloud_migrations_migration
where
  lifecycle_state = 'ACTIVE';
```

### Find migrations by display name
Explore which migration projects match a specific naming pattern to organize and track related migrations. This allows administrators to easily locate migrations that follow specific naming conventions.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration
where
  display_name like 'migration_%';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration
where
  display_name like 'migration_%';
```

### Find migrations in a specific compartment
Explore which migration projects exist within a particular compartment to better manage resources and track compartment-specific migrations. This is useful for departmental or project-based migration tracking.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration
where
  compartment_id = 'ocid1.compartment.oc1..exampleuniqueID';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration
where
  compartment_id = 'ocid1.compartment.oc1..exampleuniqueID';
```

### List completed migrations
Determine which migration projects have been successfully completed to track migration progress and verify successful transfers. This helps in understanding which workloads have completed their migration journey.

```sql+postgres
select
  id,
  display_name,
  time_created,
  time_updated
from
  oci_cloud_migrations_migration
where
  is_completed = true;
```

```sql+sqlite
select
  id,
  display_name,
  time_created,
  time_updated
from
  oci_cloud_migrations_migration
where
  is_completed = 1;
```

### Get migration details with associated replication schedule
Explore which migration projects have specific replication schedules attached to understand the data synchronization patterns. This is useful for tracking how frequently data is being replicated as part of the migration process.

```sql+postgres
select
  m.id,
  m.display_name,
  m.lifecycle_state,
  m.replication_schedule_id,
  r.display_name as schedule_name,
  r.execution_recurrences
from
  oci_cloud_migrations_migration as m
left join
  oci_cloud_migrations_replication_schedule as r
on
  m.replication_schedule_id = r.id
where
  m.replication_schedule_id is not null;
```

```sql+sqlite
select
  m.id,
  m.display_name,
  m.lifecycle_state,
  m.replication_schedule_id,
  r.display_name as schedule_name,
  r.execution_recurrences
from
  oci_cloud_migrations_migration as m
left join
  oci_cloud_migrations_replication_schedule as r
on
  m.replication_schedule_id = r.id
where
  m.replication_schedule_id is not null;
```