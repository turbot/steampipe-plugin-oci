---
title: "Steampipe Table: oci_cloud_migrations_replication_schedule - Query OCI Cloud Migrations Replication Schedules using SQL"
description: "Allows users to query OCI Cloud Migrations Replication Schedules."
---

# Table: oci_cloud_migrations_replication_schedule - Query OCI Cloud Migrations Replication Schedules using SQL

Oracle Cloud Infrastructure (OCI) Cloud Migrations service helps enterprises to migrate their on-premises or other cloud workloads to OCI. A Replication Schedule defines when and how frequently data is replicated from source environments to OCI during migration. It controls the synchronization process to minimize data loss and downtime during cutover.

## Table Usage Guide

The `oci_cloud_migrations_replication_schedule` table provides insights into replication schedules within OCI Cloud Migrations service. As a migration engineer or cloud administrator, explore schedule-specific details through this table, including recurrence patterns, state, and timing. Utilize it to monitor data synchronization frequency, verify proper replication timing, and manage data transfer patterns during migration.

## Examples

### Basic info
Explore the fundamental details of your replication schedules to understand their configuration and status. This helps in monitoring how frequently your data is being synchronized during migration to OCI.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_replication_schedule;
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_replication_schedule;
```

### List active replication schedules
Identify which replication schedules are currently active to monitor ongoing data synchronization processes. This information is crucial for ensuring that your migration maintains up-to-date data copies.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  time_created
from
  oci_cloud_migrations_replication_schedule
where
  lifecycle_state = 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  time_created
from
  oci_cloud_migrations_replication_schedule
where
  lifecycle_state = 'ACTIVE';
```

### Filter replication schedules by display name
Discover replication schedules that follow specific naming conventions to organize and categorize your data synchronization processes. This helps in managing schedules for different types of migrations or workloads.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  display_name like 'schedule_%';
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  display_name like 'schedule_%';
```

### Find replication schedules in a specific compartment
Explore which replication schedules exist within a particular compartment to organize and track compartment-specific migration processes. This is useful for departmental or project-based migration management.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  compartment_id = 'ocid1.compartment.oc1..exampleuniqueID';
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  compartment_id = 'ocid1.compartment.oc1..exampleuniqueID';
```

### Find schedules using specific recurrence patterns
Analyze replication schedules with particular recurrence patterns to understand the frequency of data synchronization. This helps in identifying schedules with daily or other specific replication timings.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  execution_recurrences like '%FREQ=DAILY%';
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  lifecycle_state
from
  oci_cloud_migrations_replication_schedule
where
  execution_recurrences like '%FREQ=DAILY%';
```

### Find recently created replication schedules
Identify replication schedules that were created in the last 30 days to monitor recent changes to your migration synchronization processes. This helps in tracking newly established data replication patterns.

```sql+postgres
select
  id,
  display_name,
  execution_recurrences,
  time_created
from
  oci_cloud_migrations_replication_schedule
where
  time_created >= current_date - interval '30' day;
```

```sql+sqlite
select
  id,
  display_name,
  execution_recurrences,
  time_created
from
  oci_cloud_migrations_replication_schedule
where
  time_created >= datetime('now', '-30 day');
```