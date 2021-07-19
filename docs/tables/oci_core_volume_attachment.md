# Table: oci_core_volume_attachment

You can attach a volume to an instance in order to expand the available storage on the instance. If you specify iSCSI as the volume attachment type, you must also connect and mount the volume from the instance for the volume to be usable.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment;
```

### List idle volume attachments

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  lifecycle_state <> 'ATTACHED';
```

### List read only volume attachments

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_attachment
where
  is_read_only;
```
