---
title: "Steampipe Table: oci_core_boot_volume_metric_write_ops - Query OCI Core Boot Volume Metrics using SQL"
description: "Allows users to query Boot Volume Write Operations Metrics in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_core_boot_volume_metric_write_ops - Query OCI Core Boot Volume Metrics using SQL

Boot Volume in OCI is a block storage volume that contains the image used to boot a Compute instance. These boot volumes are reliable and durable with built-in redundancy within the availability domain. The write operations metrics provide insights into the write operations performed on the boot volume.

## Table Usage Guide

The `oci_core_boot_volume_metric_write_ops` table provides insights into write operations metrics of boot volumes in OCI. As a system administrator or a DevOps engineer, explore details of write operations on boot volumes through this table, including the number of operations, average size, and total bytes written. Utilize it to monitor and optimize the performance of boot volumes, ensuring efficient operation of your Compute instances in OCI.

## Examples

### Basic info
Explore which boot volumes in your OCI Core have the highest activity by analyzing write operations. This can help determine where potential bottlenecks or high usage might be occurring.

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
  oci_core_boot_volume_metric_write_ops
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
  oci_core_boot_volume_metric_write_ops
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 1000 average write ops
Analyze the intervals where the average write operations exceed 1000 for boot volumes. This is useful for identifying periods of high activity and potential performance issues.

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
  oci_core_boot_volume_metric_write_ops
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
  oci_core_boot_volume_metric_write_ops
where
  average > 1000
order by
  id,
  timestamp;
```

### Intervals where volumes exceed 8000 max write ops
Explore instances where the maximum write operations on boot volumes exceed a certain threshold. This can help in identifying potential bottlenecks or performance issues in the system.

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
  oci_core_boot_volume_metric_write_ops
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
  oci_core_boot_volume_metric_write_ops
where
  maximum > 8000
order by
  id,
  timestamp;
```

### Read, Write, and Total IOPS
Gain insights into the input/output operations of your boot volume by assessing both read and write operations. This allows you to monitor and optimize the performance of your storage system over time.

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