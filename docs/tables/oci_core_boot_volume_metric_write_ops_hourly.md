---
title: "Steampipe Table: oci_core_boot_volume_metric_write_ops_hourly - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query hourly write operations metrics for OCI Core Boot Volumes."
---

# Table: oci_core_boot_volume_metric_write_ops_hourly - Query OCI Core Boot Volume Metrics using SQL

Oracle Cloud Infrastructure's Core Boot Volume is a block storage volume that contains the image used to boot a Compute instance. These boot volumes are reliable and durable, with built-in redundancy to protect your data against failure. They also offer high performance and a large storage capacity.

## Table Usage Guide

The `oci_core_boot_volume_metric_write_ops_hourly` table provides insights into the hourly write operations metrics of OCI Core Boot Volumes. As a data analyst or a cloud operations engineer, you can use this table to monitor and analyze the write operations performance of your boot volumes on an hourly basis. This can be particularly useful for identifying potential issues, optimizing performance, and ensuring the efficient use of resources.

## Examples

### Basic info
Explore the hourly write operations on boot volumes in Oracle Cloud Infrastructure. This can help assess the volume's performance over time and identify any potential issues or trends.

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
  oci_core_boot_volume_metric_write_ops_hourly
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
  oci_core_boot_volume_metric_write_ops_hourly
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Explore instances where the average write operations exceed 1000 on an hourly basis. This can be useful for identifying potential periods of high activity or system stress.

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
  oci_core_boot_volume_metric_write_ops_hourly
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
  oci_core_boot_volume_metric_write_ops_hourly
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Determine the instances where the maximum write operations on boot volumes exceed a set threshold. This can help in identifying potential performance issues and planning for capacity upgrades.

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
  oci_core_boot_volume_metric_write_ops_hourly
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
  oci_core_boot_volume_metric_write_ops_hourly
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Analyze the performance of your boot volume by observing the average, maximum, and minimum input/output operations per second (IOPS). This can help you understand how your system is performing and where potential bottlenecks might be occurring.

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