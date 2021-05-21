# Table: oci_cloud_guard_managed_list

OCI Cloud Guard Managed Lists provide a centralized location for detector rule configuration. You can define a list one time and use it in multiple rules.

## Examples

### Basic info

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_managed_list;
```

### List active managed lists

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_managed_list
where
  lifecycle_state = 'ACTIVE';
```
