# Table: oci_nosql_table_metric_read_throttle_count_hourly

OCI Monitoring Metrics provide data about the performance of your systems.  The `oci_nosql_table_metric_read_throttle_count_hourly` table provides metric statistics at 60 minute intervals for the most recent 60 days.


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
  oci_nosql_table_metric_read_throttle_count_hourly
order by
  table_name,
  timestamp;
```

### Intervals where read throttle count exceded 100 average
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
  oci_nosql_table_metric_read_throttle_count_hourly
where
  average > 100
order by
  table_name,
  timestamp;
```
