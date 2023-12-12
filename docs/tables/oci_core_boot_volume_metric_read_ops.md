---
title: "Steampipe Table: oci_core_boot_volume_metric_read_ops - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query OCI Core Boot Volume Read Operations Metrics."
---

# Table: oci_core_boot_volume_metric_read_ops - Query OCI Core Boot Volume Metrics using SQL

Oracle Cloud Infrastructure's Core Boot Volume is a block storage volume that contains the image used to boot a Compute instance. These Boot Volumes are reliable, high-performance storage volumes that can persistently store data and can be used as boot volumes for instances. They offer consistent, low-latency performance, and a variety of management features, such as backup and restore operations, cloning, and volume expansion.

## Table Usage Guide

The `oci_core_boot_volume_metric_read_ops` table provides insights into the read operations metrics of Boot Volumes within Oracle Cloud Infrastructure's Core service. As a system administrator or DevOps engineer, explore metric-specific details through this table, including the volume ID, namespace, metric timestamp, and read operation statistics. Utilize it to monitor and analyze the performance of your Boot Volumes, identify any unusual read operation patterns, and ensure optimal performance of your instances.

## Examples

### Basic info
Explore the performance metrics of boot volumes in your Oracle Cloud Infrastructure environment. This query helps in identifying any unusual activity or performance degradation by analyzing metrics such as minimum, maximum, average, and total read operations over time.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Analyze the settings to understand periods where the average read operations surpass a threshold of 1000. This could be useful in identifying potential system overloads or performance bottlenecks.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
where
  average > 1000
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max read ops
Explore instances where the read operations on boot volumes exceed a certain threshold to monitor system performance and identify potential bottlenecks. This can be useful in optimizing your system configuration for improved efficiency.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
where
  maximum > 8000
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_core_boot_volume_metric_read_ops
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Explore the performance of your system by analyzing input/output operations. This query helps to understand the average, maximum, and minimum read/write operations over time, assisting in identifying potential issues or areas for improvement.

```sql+postgres
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

```sql+sqlite
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