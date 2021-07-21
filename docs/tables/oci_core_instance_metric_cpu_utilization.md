# Table: oci_core_instance_metric_cpu_utilization

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_core_instance_metric_cpu_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  oci_core_instance_metric_cpu_utilization
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
  oci_core_instance_metric_cpu_utilization
where average > 80
order by
  id,
  timestamp;
```
