# Table: oci_nosql_table_metric_storage_utilization_hourly

OCI Monitoring metrics provide data about the performance of your systems. The `oci_nosql_table_metric_storage_utilization_hourly` table provides metric statistics at 60 minute intervals for the most recent 60 days.

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
  oci_nosql_table_metric_storage_utilization_hourly
order by
  name,
  timestamp;
```

### Intervals where storage utilization exceeded 20GB average

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
  oci_nosql_table_metric_storage_utilization_hourly
where
  average > 20 
order by
  name,
  timestamp;
```
