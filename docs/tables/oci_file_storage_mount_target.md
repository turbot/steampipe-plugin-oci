# Table: oci_file_storage_mount_target

A mount target is an NFS endpoint that lives in a VCN subnet of your choice and provides network access for file systems. The mount target provides the IP address or DNS name that is used together with a unique export path to mount the file system.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  availability_domain,
  time_created
from
  oci_file_storage_mount_target;
```


## List mount targets that are not active

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_mount_target
where
  lifecycle_state <> 'ACTIVE';
```
