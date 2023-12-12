---
title: "Steampipe Table: oci_core_boot_volume_metric_read_ops_hourly - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query OCI Core Boot Volume Read Operations Metrics on an hourly basis."
---

# Table: oci_core_boot_volume_metric_read_ops_hourly - Query OCI Core Boot Volume Metrics using SQL

Oracle Cloud Infrastructure (OCI) Core Boot Volume is a persistent, block-level storage volume that you can attach to a single instance. The boot volume contains the image of the operating system running on your instance. The `oci_core_boot_volume_metric_read_ops_hourly` table provides data related to the read operations performed on the boot volumes, aggregated on an hourly basis.

## Table Usage Guide

The `oci_core_boot_volume_metric_read_ops_hourly` table provides insights into the read operations metrics of OCI Core Boot Volumes. As a cloud engineer or system administrator, you can use this table to monitor and analyze the read operations on boot volumes, which can be crucial for performance tuning and troubleshooting. This table can be particularly useful in identifying volumes with high read operations, which might indicate a need for capacity planning or performance optimization.

## Examples

### Basic info
Analyze the settings to understand the performance of boot volumes in Oracle Cloud Infrastructure over time. This query can be used to monitor the read operations, allowing you to pinpoint any unusual activity or potential bottlenecks.

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
  oci_core_boot_volume_metric_read_ops_hourly
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
  oci_core_boot_volume_metric_read_ops_hourly
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average read ops
Identify instances where the average read operations on boot volumes surpass 1000 within an hour. This can help pinpoint potential areas of high workload and facilitate proactive system management.

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
  oci_core_boot_volume_metric_read_ops_hourly
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
  oci_core_boot_volume_metric_read_ops_hourly
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max read ops
Assess the instances where the maximum read operations on boot volumes exceed a threshold of 8000. This can be beneficial in identifying periods of high traffic or potential performance issues within your server infrastructure.

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
  oci_core_boot_volume_metric_read_ops_hourly
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
  oci_core_boot_volume_metric_read_ops_hourly
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Determine the areas in which input/output operations per second (IOPS) are occurring, providing a comprehensive view of both read and write operations. This can help optimize system performance by identifying potential bottlenecks or areas for improvement.

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
  oci_core_boot_volume_metric_read_ops_hourly as r,
  oci_core_boot_volume_metric_write_ops_hourly as w
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
  oci_core_boot_volume_metric_read_ops_hourly as r,
  oci_core_boot_volume_metric_write_ops_hourly as w
where 
  r.id = w.id
  and r.timestamp = w.timestamp
order by
  r.id,
  r.timestamp;
```