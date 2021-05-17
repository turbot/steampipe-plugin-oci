# Table: oci_logging_log

The Oracle Cloud Infrastructure Logging service is a highly scalable and fully managed single pane of glass for all the logs in your tenancy. Logging provides access to logs from Oracle Cloud Infrastructure resources. These logs include critical diagnostic information that describes how resources are performing and being accessed.

## Examples

### Basic info

```sql
select
  id,
  log_group_id,
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
  oci_logging_log
where
  lifecycle_state = 'INACTIVE';
```

### List VCN Subnets where flow log is enabled

```sql
select
  configuration -> 'source' ->> 'resource' as subnet_id,
  configuration -> 'source' ->> 'service' as service,
  lifecycle_state
from
  oci_logging_log
where
  configuration -> 'source' ->> 'service' = 'flowlogs'
  and lifecycle_state = 'ACTIVE';
```
