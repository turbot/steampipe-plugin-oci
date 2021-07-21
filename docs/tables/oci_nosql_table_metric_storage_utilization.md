# Table: oci_nosql_table_metric_storage_utilization

OCI Monitoring metrics provide data about the performance of your systems. The `oci_nosql_table_metric_storage_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  oci_nosql_table_metric_storage_utilization
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
  oci_nosql_table_metric_storage_utilization
where
  average > 20 
order by
  name,
  timestamp;
```
