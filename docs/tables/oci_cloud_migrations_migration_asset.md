---
title: "Steampipe Table: oci_cloud_migrations_migration_asset - Query OCI Cloud Migrations Migration Assets using SQL"
description: "Allows users to query OCI Cloud Migrations Migration Assets."
---

# Table: oci_cloud_migrations_migration_asset - Query OCI Cloud Migrations Migration Assets using SQL

Oracle Cloud Infrastructure (OCI) Cloud Migrations service helps enterprises to migrate their on-premises or other cloud workloads to OCI. A Migration Asset represents a specific resource being migrated, such as a virtual machine, database, or application component, and contains metadata about the source and target configurations of that resource.

## Table Usage Guide

The `oci_cloud_migrations_migration_asset` table provides insights into migration assets within OCI Cloud Migrations service. As a migration engineer or cloud administrator, explore asset-specific details through this table, including their type, state, and replication details. Utilize it to inventory migration resources, track asset replication status, and manage dependencies between migrated resources.

## Examples

### Basic info
Explore the fundamental details of assets being migrated to Oracle Cloud to gain an overview of your migration landscape. This can help in tracking the progress and status of individual components within your migration project.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  type,
  lifecycle_state
from
  oci_cloud_migrations_migration_asset;
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  type,
  lifecycle_state
from
  oci_cloud_migrations_migration_asset;
```

### Find assets for a specific migration
Analyze which assets are associated with a particular migration project to track its scope and progress. This helps in understanding the composition and status of resources within a specific migration.

```sql+postgres
select
  id,
  display_name,
  type,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_asset
where
  migration_id = 'ocid1.ocmmigration.oc1.region.example';
```

```sql+sqlite
select
  id,
  display_name,
  type,
  lifecycle_state,
  time_created
from
  oci_cloud_migrations_migration_asset
where
  migration_id = 'ocid1.ocmmigration.oc1.region.example';
```

### List assets by lifecycle state
Identify migration assets in a specific state to monitor progress and troubleshoot issues. This query helps in tracking which components are actively being migrated and which may require attention.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  type
from
  oci_cloud_migrations_migration_asset
where
  lifecycle_state = 'ACTIVE';
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  type
from
  oci_cloud_migrations_migration_asset
where
  lifecycle_state = 'ACTIVE';
```

### Find assets of a specific type
Discover the details of a particular type of resource being migrated to better manage and track similar assets. This can help in planning and coordinating the migration of specific workload components, such as virtual machines.

```sql+postgres
select
  id,
  migration_id,
  display_name,
  lifecycle_state
from
  oci_cloud_migrations_migration_asset
where
  type = 'VM';
```

```sql+sqlite
select
  id,
  migration_id,
  display_name,
  lifecycle_state
from
  oci_cloud_migrations_migration_asset
where
  type = 'VM';
```

### Find assets with snapshots
Examine which migration assets have snapshots configured to understand their replication and recovery capabilities. This helps in verifying data protection mechanisms during the migration process.

```sql+postgres
select
  id,
  display_name,
  type,
  lifecycle_state,
  snapshots
from
  oci_cloud_migrations_migration_asset
where
  snapshots is not null;
```

```sql+sqlite
select
  id,
  display_name,
  type,
  lifecycle_state,
  snapshots
from
  oci_cloud_migrations_migration_asset
where
  snapshots is not null;
```

### Analyze asset dependencies
Identify the dependencies between migration assets to understand the relationship structure and potential migration sequence. This is crucial for planning a proper migration order and preventing dependency-related issues.

```sql+postgres
select
  id,
  display_name,
  type,
  depends_on,
  depended_on_by
from
  oci_cloud_migrations_migration_asset
where
  depends_on is not null or depended_on_by is not null;
```

```sql+sqlite
select
  id,
  display_name,
  type,
  depends_on,
  depended_on_by
from
  oci_cloud_migrations_migration_asset
where
  depends_on is not null or depended_on_by is not null;
```