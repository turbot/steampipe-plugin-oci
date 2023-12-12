---
title: "Steampipe Table: oci_core_volume_backup - Query OCI Core Volume Backups using SQL"
description: "Allows users to query OCI Core Volume Backups."
---

# Table: oci_core_volume_backup - Query OCI Core Volume Backups using SQL

Oracle Cloud Infrastructure (OCI) Core Volume Backup is a service that provides automated backups of block volumes. The backups are crash-consistent and can be used to create new volumes or restore existing volumes. This service ensures data durability and protection for your OCI resources.

## Table Usage Guide

The `oci_core_volume_backup` table provides insights into volume backups within OCI Core. As a system administrator, explore backup-specific details through this table, including backup size, status, and associated metadata. Utilize it to uncover information about backups, such as those that are available, the ones that are in progress, and the verification of backup policies.

## Examples

### Basic info
Explore which backups of your core volumes in Oracle Cloud Infrastructure have been created, along with their source type and current lifecycle state. This can help you manage and track your backups, ensuring data safety and availability.

```sql+postgres
select
  id,
  display_name,
  volume_id,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_volume_backup;
```

```sql+sqlite
select
  id,
  display_name,
  volume_id,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_volume_backup;
```


### List manual volume backup
Explore which volume backups have been manually created. This can be useful for auditing purposes, ensuring that backups are being created as intended and identifying any potential issues with backup procedures.

```sql+postgres
select
  id,
  display_name,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_core_volume_backup
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
  oci_core_volume_backup
where
  source_type = 'MANUAL';
```