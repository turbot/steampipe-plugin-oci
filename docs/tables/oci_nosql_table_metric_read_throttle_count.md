# Table: oci_nosql_table_metric_read_throttle_count

OCI Monitoring metrics provide data about the performance of your systems. The `oci_nosql_table_metric_read_throttle_count` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_read_throttle_count
order by
  name,
  timestamp;
```

### Intervals where read throttle count exceeded 100 average

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_read_throttle_count
where
  average > 100
order by
  name,
  timestamp;
```
