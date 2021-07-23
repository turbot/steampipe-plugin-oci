# Table: oci_mysql_db_system_metric_memory_utilization

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_mysql_db_system_metric_memory_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization
order by
  id,
  timestamp;
```

### Memory Utilization Over 80% average

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_memory,
  round(maximum::numeric,2) as max_memory,
  round(average::numeric,2) as avg_memory,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization
where average > 80
order by
  id,
  timestamp;
```
