---
title: "Steampipe Table: oci_core_boot_volume_metric_read_ops_daily - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query daily read operations metrics for OCI Core Boot Volumes."
---

# Table: oci_core_boot_volume_metric_read_ops_daily - Query OCI Core Boot Volume Metrics using SQL

The Oracle Cloud Infrastructure (OCI) Core Boot Volume is a persistent and durable block storage volume in OCI. It provides high performance, reliability, and scalability for your applications. Boot volumes are used to boot up an instance and contain the image of the operating system running on your instance.

## Table Usage Guide

The `oci_core_boot_volume_metric_read_ops_daily` table provides insights into the daily read operations metrics for OCI Core Boot Volumes. As a Database Administrator or a DevOps engineer, you can use this table to monitor the performance of your boot volumes, which can help in analyzing the workload and making data-driven decisions for optimizing resource allocation. Utilize it to uncover information about the read operations, such as the volume of data read from your boot volumes, and the time taken for these operations, which can be crucial for performance tuning and troubleshooting.

## Examples

### Basic info
Analyze the daily read operations of boot volumes in Oracle Cloud Infrastructure (OCI) to understand the range, average, total count, and specific instances of these operations. This could help in monitoring the performance and usage patterns of the boot volumes.

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
  oci_core_boot_volume_metric_read_ops_daily
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
  oci_core_boot_volume_metric_read_ops_daily
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Analyze the intervals where the average read operations on boot volumes exceed a thousand. This analysis can help identify periods of high demand or potential performance issues.

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
  oci_core_boot_volume_metric_read_ops_daily
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
  oci_core_boot_volume_metric_read_ops_daily
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max read ops
Determine the instances where the maximum read operations on a boot volume surpass 8000, allowing you to identify potential bottlenecks or high-demand periods in your system. This can be crucial for capacity planning and optimizing system performance.

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
  oci_core_boot_volume_metric_read_ops_daily
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
  oci_core_boot_volume_metric_read_ops_daily
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Determine the areas in which the input/output operations per second (IOPS) for a boot volume are at their maximum, minimum, and average. This analysis can help optimize the performance of your system by identifying potential bottlenecks or areas for improvement.

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
  oci_core_boot_volume_metric_read_ops_daily as r,
  oci_core_boot_volume_metric_write_ops_daily as w
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
  oci_core_boot_volume_metric_read_ops_daily as r,
  oci_core_boot_volume_metric_write_ops_daily as w
where 
  r.id = w.id
  and r.timestamp = w.timestamp
order by
  r.id,
  r.timestamp;
```