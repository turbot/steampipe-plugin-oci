---
title: "Steampipe Table: oci_core_boot_volume_metric_write_ops_daily - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query daily write operations metrics of OCI Core Boot Volumes."
---

# Table: oci_core_boot_volume_metric_write_ops_daily - Query OCI Core Boot Volume Metrics using SQL

OCI Core Boot Volumes are persistent storage devices that provide durable block storage for instances within Oracle Cloud Infrastructure (OCI). They are used as the primary storage device for hosting an instance's operating system, system software, and other boot volume data. Boot Volumes offer consistent, low-latency performance and are integrated with OCI's security and management policies.

## Table Usage Guide

The `oci_core_boot_volume_metric_write_ops_daily` table provides insights into the daily write operations of OCI Core Boot Volumes. As a cloud engineer, you can use this table to monitor and analyze the write performance of your boot volumes to optimize your resource usage and troubleshoot issues. This table can be particularly useful for identifying high-utilization periods and potential bottlenecks in your system.

## Examples

### Basic info
Analyze the daily write operations of boot volumes in your OCI Core infrastructure. This query helps you understand the performance trends and potential bottlenecks in your system by examining metrics such as minimum, maximum, average operations, and total operations.

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
  oci_core_boot_volume_metric_write_ops_daily
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
  oci_core_boot_volume_metric_write_ops_daily
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Explore intervals where the average daily write operations on your boot volumes exceed 1000. This can help identify potential performance issues or unusual activity in your system.

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
  oci_core_boot_volume_metric_write_ops_daily
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
  oci_core_boot_volume_metric_write_ops_daily
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Determine the instances where the maximum daily write operations on boot volumes exceed 8000, enabling you to identify potential areas of high disk activity and optimize system performance.

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
  oci_core_boot_volume_metric_write_ops_daily
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
  oci_core_boot_volume_metric_write_ops_daily
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Determine the areas in which the average, maximum, and minimum Input/Output operations per second (IOPS) are analyzed for both read and write operations. This can be useful to understand the performance and efficiency of your boot volumes over time.

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