# Table: oci_database_autonomous_database_metric_storage_utilization_daily

OCI Monitoring metrics explorer provide data about the performance of your systems.  The `oci_database_autonomous_database_metric_storage_utilization_daily` table provides metric statistics at 24 hour intervals for the most recent 90 days.

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
  oci_database_autonomous_database_metric_storage_utilization_daily
order by
  id,
  timestamp;
```

### Storage Utilization Over 80% average

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_storage,
  round(maximum::numeric,2) as max_storage,
  round(average::numeric,2) as avg_storage,
  sample_count
from
  oci_database_autonomous_database_metric_storage_utilization_daily
where average > 80
order by
  id,
  timestamp;
```