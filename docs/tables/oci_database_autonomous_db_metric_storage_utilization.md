# Table: oci_database_autonomous_database_metric_storage_utilization

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_database_autonomous_db_metric_storage_utilization` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  oci_database_autonomous_database_metric_storage_utilization
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
  oci_database_autonomous_database_metric_storage_utilization
where average > 80
order by
  id,
  timestamp;
```
