# Table: oci_core_boot_volume_attachment

OCI core boot volume attachment provides the boot volume attachments details.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume_attachment;
```

### List volume attachments witch are not attached

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume_attachment
where
  lifecycle_state <> 'ATTACHED';
```
