# Table: oci_nosql_table_metric_storage_utilization_hourly

OCI Monitoring metrics provide data about the performance of your systems. The `oci_nosql_table_metric_storage_utilization_hourly` table provides metric statistics at 60 minute intervals for the most recent 60 days.

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
  oci_nosql_table_metric_storage_utilization_hourly
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
  oci_nosql_table_metric_storage_utilization_hourly
where
  average > 20 
order by
  table_name,
  timestamp;
```
