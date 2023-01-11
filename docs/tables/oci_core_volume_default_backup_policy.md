# Table: oci_core_volume_default_backup_policy

A policy for automatically creating volume backups according to a recurring schedule. Has a set of one or more schedules that control when and how backups are created.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  time_created,
  region,
  tags
from
  oci_core_volume_default_backup_policy;
```

### Get the destination region for each volume backup policy

```sql
select
  id,
  display_name,
  destination_region
from
  oci_core_volume_default_backup_policy;
```

### Get schedule info for each volume backup policy

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
  oci_core_volume_default_backup_policy,
  jsonb_array_elements(schedules) as s;
```

### List volume back policies that create full backups

```sql
select
  id,
  display_name,
  s ->> 'backupType' as backup_type
from
  oci_core_volume_default_backup_policy,
  jsonb_array_elements(schedules) as s
where
  s ->> 'backupType' = 'FULL';
```