# Table: oci_core_boot_volume_metric_write_ops

OCI Monitoring metrics explorer provide data about the performance of your systems. The `oci_core_boot_volume_metric_write_ops` table provides metric statistics at 5 minute intervals for the most recent 5 days.

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
  oci_core_boot_volume_metric_write_ops
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops

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
  oci_core_boot_volume_metric_write_ops
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops

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
  oci_core_boot_volume_metric_write_ops
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
  oci_core_boot_volume_metric_read_ops as r,
  oci_core_boot_volume_metric_write_ops as w
where 
  r.id = w.id
  and r.timestamp = w.timestamp
order by
  r.id,
  r.timestamp;
```
