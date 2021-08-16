# Table: oci_nosql_table_metric_write_throttle_count_daily

OCI Monitoring Metrics provide data about the performance of your systems. The `oci_nosql_table_metric_write_throttle_count_daily` table provides metric statistics at 24 hour intervals for the most recent 90 days.

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
  oci_nosql_table_metric_write_throttle_count_daily
order by
  name,
  timestamp;
```

### Intervals where write throttle count exceeded 100 average

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
  oci_nosql_table_metric_write_throttle_count_daily
where
  average > 100
order by
  name,
  timestamp;
```
