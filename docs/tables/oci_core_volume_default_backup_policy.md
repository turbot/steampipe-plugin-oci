---
title: "Steampipe Table: oci_core_volume_default_backup_policy - Query OCI Core Volume Default Backup Policies using SQL"
description: "Allows users to query OCI Core Volume Default Backup Policies"
---

# Table: oci_core_volume_default_backup_policy - Query OCI Core Volume Default Backup Policies using SQL

The OCI Core Volume Default Backup Policy is a resource within Oracle Cloud Infrastructure (OCI) that provides automated backups of block volume data. The policy determines the frequency of automatic backups and the retention period for these backups. It is a crucial component of the OCI Block Volume service, helping to ensure data durability and protection.

## Table Usage Guide

The `oci_core_volume_default_backup_policy` table provides insights into the default backup policies associated with OCI Core Volumes. As a database administrator or DevOps engineer, use this table to explore policy-specific details, such as backup frequency and retention periods. This can be particularly useful for maintaining regular data backups, ensuring data durability, and planning disaster recovery strategies.

## Examples

### Basic info
Explore which default backup policies have been created and when, to gain insights into the history and management of your data backups. This can help in assessing the regularity and effectiveness of your backup strategies.

```sql
select
  id,
  display_name,
  time_created
from
  oci_core_volume_default_backup_policy;
```

### Get schedule info for each volume backup policy
This query is useful for gaining insights into the scheduling of volume backup policies. It helps in understanding the timing, frequency, and retention period of backups, which can assist in optimizing storage management and disaster recovery plans.

```sql
select
  id,
  display_name,
  s ->> 'backupType' as backup_type,
  s ->> 'dayOfMonth' as day_of_month,
  s ->> 'hourOfDay' as hour_of_day,
  s ->> 'offsetSeconds' as offset_seconds,
  s ->> 'period' as period,
  s ->> 'retentionSeconds' as retention_seconds,
  s ->> 'timeZone' as time_zone
from
  oci_core_volume_default_backup_policy,
  jsonb_array_elements(schedules) as s;
```