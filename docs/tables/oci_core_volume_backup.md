# Table: oci_core_volume_backup

The backups feature of the Oracle Cloud Infrastructure Block Volume service lets you make a point-in-time snapshot of the data on a block volume

## Examples

### Basic info

```sql
select
  id,
  display_name,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_new.oci_core_volume_backup;
```

### List manual volume backup

```sql
select
  id,
  display_name,
  source_type,
  time_created,
  type,
  lifecycle_state as state
from
  oci_new.oci_core_volume_backup
where
  source_type = 'MANUAL';
```
