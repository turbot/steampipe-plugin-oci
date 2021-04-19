# Table: oci_logging_log_group

The Oracle Cloud Infrastructure Logging service is a highly scalable and fully managed single pane of glass for all the logs in your tenancy. Logging provides access to logs from Oracle Cloud Infrastructure resources. These logs include critical diagnostic information that describes how resources are performing and being accessed.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_logging_log;
```

### List logs which are inactive

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created
from
  oci_logging_log_
where
  lifecycle_state = 'INACTIVE';
```
