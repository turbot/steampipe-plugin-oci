---
title: "Steampipe Table: oci_core_volume_backup_policy - Query OCI Core Volume Backup Policies using SQL"
description: "Allows users to query OCI Core Volume Backup Policies."
---

# Table: oci_core_volume_backup_policy - Query OCI Core Volume Backup Policies using SQL

A Core Volume Backup Policy in Oracle Cloud Infrastructure (OCI) is a resource that defines a set of rules for automatic backups of block volumes. These policies are designed to automate the process of backing up your data, ensuring that it's protected and can be restored if necessary. They offer flexibility in terms of scheduling, allowing you to specify when backups should occur and how long they should be retained.

## Table Usage Guide

The `oci_core_volume_backup_policy` table provides insights into the backup policies of block volumes within OCI. As a system administrator, you can use this table to explore policy-specific details, including backup schedules, retention periods, and associated metadata. This can be particularly useful for ensuring data protection and compliance with organizational backup policies.

## Examples

### Basic info
Explore which volume backup policies have been created in specific regions, along with their creation times and any applied tags. This can be useful for managing and tracking backup policies across different geographies.

```sql
select
  id,
  display_name,
  time_created,
  region,
  tags
from
  oci_core_volume_backup_policy;
```


### Get the destination region for each volume backup policy
Determine the destination region for each volume backup policy to help manage and optimize data storage across different regions.

```sql
select
  id,
  display_name,
  destination_region
from
  oci_core_volume_backup_policy;
```


### Get schedule info for each volume backup policy
Explore the scheduling details for each volume backup policy. This can help in identifying when and how frequently backups are taken, thus enabling efficient resource planning and data recovery strategies.

```sql
select
  id,
  display_name,
  s ->> 'backupType' as backup_type,
  s ->> 'dayOfMonth' as day_of_month,
  s ->> 'dayOfWeek' as day_of_week,
  s ->> 'hourOfDay' as hour_of_day,
  s ->> 'month' as month,
  s ->> 'offsetSeconds' as offset_econds,
  s ->> 'offsetType' as offset_type,
  s ->> 'period' as offset_econds,
  s ->> 'retentionSeconds' as retention_seconds,
  s ->> 'timeZone' as time_zone
from
  oci_core_volume_backup_policy,
  jsonb_array_elements(schedules) as s;
```


### List volume back policies that create full backups
Explore which volume backup policies are set to create full backups. This can be beneficial to ensure important data is completely backed up and to identify any areas that may require a change in backup strategy.

```sql
select
  id,
  display_name,
  s ->> 'backupType' as backup_type
from
  oci_core_volume_backup_policy,
  jsonb_array_elements(schedules) as s
where
  s ->> 'backupType' = 'FULL';
```