# Table: oci_nosql_table_metric_write_throttle_count

OCI Monitoring Metrics provide data about the performance of your systems.  The `oci_nosql_table_metric_write_throttle_count` table provides metric statistics at 5 minute intervals for the most recent 5 days.


## Examples

### Basic info

```sql
select
  table_name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count
order by
  table_name,
  timestamp;
```

### Intervals where write throttle count exceded 100 average
```sql
select
  table_name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count
where
  average > 100
order by
  table_name,
  timestamp;
```
