# Table: oci_core_boot_volume_metric_read_ops_daily

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_core_boot_volume_metric_read_ops_daily` table provides metric statistics at 24 hour intervals for the most recent 90 days.

## Examples

### Basic info

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops_daily
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops_daily
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max read ops

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops_daily
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS

```sql
select 
  r.id,
  r.timestamp,
  round(r.average) + round(w.average) as iops_avg,
  round(r.average) as read_ops_avg,
  round(w.average) as write_ops_avg,
  round(r.maximum) + round(w.maximum) as iops_max,
  round(r.maximum) as read_ops_max,
  round(w.maximum) as write_ops_max,
  round(r.minimum) + round(w.minimum) as iops_min,
  round(r.minimum) as read_ops_min,
  round(w.minimum) as write_ops_min
from 
  oci_core_boot_volume_metric_read_ops_daily as r,
  oci_core_boot_volume_metric_write_ops_daily as w
where 
  r.id = w.id
  and r.timestamp = w.timestamp
order by
  r.id,
  r.timestamp;
```
