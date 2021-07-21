# Table: oci_core_boot_volume_replica

The boot Volume service provides you with the capability to perform ongoing automatic asynchronous replication of boot volumes and boot volumes to other regions.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  boot_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_boot_volume_replica;
```

### List volume replicas which are not available

```sql
select
  id,
  display_name,
  boot_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_boot_volume_replica
where
  lifecycle_state <> 'AVAILABLE';
```
