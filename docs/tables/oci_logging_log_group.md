# Table: oci_logging_log_group

The Oracle Cloud Infrastructure Log groups are logical containers for organizing logs. Logs must always be inside log groups. You must create a log group to enable a log.

## Examples

### Basic info

```sql
select
  id as log_group_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_logging_log_group;
```


### List inactive log groups

```sql
select
  id as log_group_id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_logging_log_group
where
  lifecycle_state = 'INACTIVE';
```
