# Table: oci_file_storage_snapshot

The File Storage service supports snapshots for data protection of file system. Snapshots are a consistent, point-in-time view of file systems. Snapshots are copy-on-write, and scoped to the entire file system.

## Examples

### Basic info

```sql
select
  name,
  id,
  lifecycle_state as state,
  time_created,
  provenance_id,
  region
from
  oci_file_storage_snapshot;
```


## Count of snapshots created per file system

```sql
select
  file_system_id,
  count(*) as snapshots_count
from
  oci_file_storage_snapshot
group by
  file_system_id;
```


## List cloned snapshots

```sql
select
  name,
  id,
  is_clone_source
from
  oci_file_storage_snapshot
where
  is_clone_source;
```


