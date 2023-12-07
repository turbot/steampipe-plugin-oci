---
title: "Steampipe Table: oci_core_boot_volume_backup - Query OCI Core Boot Volume Backups using SQL"
description: "Allows users to query Boot Volume Backups in OCI Core service."
---

# Table: oci_core_boot_volume_backup - Query OCI Core Boot Volume Backups using SQL

A Boot Volume Backup in OCI Core service is a point-in-time, crash-consistent snapshot of a boot volume. Boot volume backups are automatically replicated across multiple availability domains for redundancy. They can be used to create new boot volumes or restore a boot volume to a specific point in time.

## Table Usage Guide

The `oci_core_boot_volume_backup` table provides insights into Boot Volume Backups within OCI Core service. As a Cloud Architect or Database Administrator, explore backup-specific details through this table, including backup state, type, and associated metadata. Utilize it to uncover information about backups, such as those in available or faulty state, the size of each backup, and the time of creation or deletion for each backup.

## Examples

### Basic info
Explore the basic details of your boot volume backups in Oracle Cloud Infrastructure. This query can be useful to understand the state and type of each backup, when it was created, and its source, providing a comprehensive overview for effective management and oversight.

```sql+postgres
select
  id,
  display_name,
  boot_volume_id,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_boot_volume_backup;
```

```sql+sqlite
select
  id,
  display_name,
  boot_volume_id,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_boot_volume_backup;
```

### List manual boot volume backups
Uncover the details of manual boot volume backups. This query helps you in identifying those backups that were manually created, providing insights into the time and state of their creation.

```sql+postgres
select
  id,
  display_name,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_boot_volume_backup
where
  source_type = 'MANUAL';
```

```sql+sqlite
select
  id,
  display_name,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_boot_volume_backup
where
  source_type = 'MANUAL';
```