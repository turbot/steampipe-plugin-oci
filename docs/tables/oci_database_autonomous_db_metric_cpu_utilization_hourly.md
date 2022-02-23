# Table: oci_database_autonomous_db_metric_cpu_utilization_hourly

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_database_autonomous_db_metric_cpu_utilization_hourly` table provides metric statistics at 1 hour intervals for the most recent 60 days.

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
  oci_database_autonomous_db_metric_cpu_utilization_hourly
order by
  id,
  timestamp;
```

### CPU Over 80% average

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization_hourly
where average > 80
order by
  id,
  timestamp;
```
