# Table: oci_core_volume_default_backup_policy

The Oracle Cloud Infrastructure Block Volume service provides you with the capability to perform volume backups and volume group backups automatically on a schedule and retain them based on the selected backup policy.

There are three Oracle defined backup policies: `Bronze`, `Silver`, and `Gold`. Each backup policy is comprised of schedules with a set backup frequency and a retention period that you cannot modify. Oracle defined backup policies are not supported for scheduled volume group backups.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  time_created
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
  s ->> 'hourOfDay' as hour_of_day,
  s ->> 'offsetSeconds' as offset_seconds,
  s ->> 'period' as period,
  s ->> 'retentionSeconds' as retention_seconds,
  s ->> 'timeZone' as time_zone
from
  oci_core_volume_default_backup_policy,
  jsonb_array_elements(schedules) as s;
```