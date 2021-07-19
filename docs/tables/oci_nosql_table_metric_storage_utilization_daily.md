# Table: oci_nosql_table_metric_storage_utilization_daily

OCI Monitoring Metrics provide data about the performance of your systems.  The `oci_nosql_table_metric_storage_utilization_daily` table provides metric statistics at 24 hour intervals for the most recent 90 days.


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
  oci_nosql_table_metric_storage_utilization_daily
order by
  table_name,
  timestamp;
```

### Intervals where storage utilization exceded 20GB average
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
  oci_nosql_table_metric_storage_utilization_daily
where
  average > 20 
order by
  table_name,
  timestamp;
```
