---
title: "Steampipe Table: oci_cloud_migrations_migration_plan - Query OCI Cloud Migrations Migration Plans using SQL"
description: "Allows users to query OCI Cloud Migrations Migration Plans."
---

# Table: oci_cloud_migrations_migration_plan - Query OCI Cloud Migrations Migration Plans using SQL

Oracle Cloud Infrastructure (OCI) Cloud Migrations service helps enterprises to migrate their on-premises or other cloud workloads to OCI. A Migration Plan defines the target environment configurations, migration strategies, and resource mappings for migrating workloads to OCI. It serves as a blueprint for how resources will be transformed and deployed during migration.

## Table Usage Guide

The `oci_cloud_migrations_migration_plan` table provides insights into migration plans within OCI Cloud Migrations service. As a migration architect or cloud administrator, explore plan-specific details through this table, including target environments, migration strategies, and resource limits. Utilize it to design and refine migration approaches, analyze resource requirements, and monitor migration planning progress.

## Examples

### Basic info
Explore the general details of your migration plans to understand their current state and when they were created. This information can help in tracking the progress and timeline of migration planning.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_plan;
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_plan;
```

### Find plans for a specific migration
Analyze which migration plans are associated with a particular migration project to track planning progress and available options. This helps in understanding the different approaches being considered for a specific migration.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_plan
where
  migration_id = 'ocid1.ocmmigration.oc1.region.example';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_plan
where
  migration_id = 'ocid1.ocmmigration.oc1.region.example';
```

### List active migration plans
Identify migration plans that are currently active to track ongoing planning efforts and focus on relevant migration strategies. This helps in prioritizing which plans should be reviewed or implemented.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  time_created
from
  oci_cloud_migrations_migration_plan
where
  lifecycle_state = 'ACTIVE';
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  time_created
from
  oci_cloud_migrations_migration_plan
where
  lifecycle_state = 'ACTIVE';
```

### Filter migration plans by display name
Discover migration plans that follow specific naming conventions to organize and locate related planning documents. This is useful for standardizing and categorizing different types of migration approaches.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  lifecycle_state
from
  oci_cloud_migrations_migration_plan
where
  display_name like 'plan_%';
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  lifecycle_state
from
  oci_cloud_migrations_migration_plan
where
  display_name like 'plan_%';
```

### Analyze migration strategies in plans
Examine the migration strategies defined in plans to understand the approach and methods being utilized for different resources. This information is crucial for verifying that appropriate migration techniques are being applied.

```sql+postgres
select
  id,
  display_name,
  strategies
from
  oci_cloud_migrations_migration_plan
where
  strategies is not null;
```

```sql+sqlite
select
  id,
  display_name,
  strategies
from
  oci_cloud_migrations_migration_plan
where
  strategies is not null;
```

### Examine target environments in migration plans
Investigate the target environments specified in migration plans to verify proper configuration of destination resources. This helps ensure that migrations are directed to the appropriate OCI environments.

```sql+postgres
select
  id,
  display_name,
  target_environments
from
  oci_cloud_migrations_migration_plan
where
  target_environments is not null;
```

```sql+sqlite
select
  id,
  display_name,
  target_environments
from
  oci_cloud_migrations_migration_plan
where
  target_environments is not null;
```
