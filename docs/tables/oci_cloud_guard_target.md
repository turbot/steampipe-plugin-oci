# Table: oci_cloud_guard_target

Target defines the scope of what Cloud Guard checks. All compartments within a target are checked in the same way and you have the same options for processing problems that are detected.

## Examples

### Basic info

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_target;
```

### List targets which are not active

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_target
where
  lifecycle_state <> 'ACTIVE';
```
