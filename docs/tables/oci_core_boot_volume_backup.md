# Table: oci_core_boot_volume_backup

Boot volume backups are point-in-time snapshots of a boot volume without application interruption or downtime.

## Examples

### Basic info

```sql
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

```sql
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
