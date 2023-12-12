---
title: "Steampipe Table: oci_core_instance_metric_cpu_utilization_daily - Query OCI Core Instance Metrics using SQL"
description: "Allows users to query daily CPU utilization metrics for OCI Core Instances."
---

# Table: oci_core_instance_metric_cpu_utilization_daily - Query OCI Core Instance Metrics using SQL

OCI Core Instances are virtual servers in the Oracle Cloud Infrastructure that offer flexible and scalable computing capabilities. They are part of the Compute service and can run both Windows and Linux operating systems. Instances are the building blocks of applications deployed in the cloud.

## Table Usage Guide

The `oci_core_instance_metric_cpu_utilization_daily` table provides insights into the daily CPU utilization metrics of OCI Core Instances. As a system administrator or a DevOps engineer, you can explore CPU usage details through this table, including maximum, minimum, and average utilization. Utilize it to monitor CPU performance, identify instances with high CPU usage, and plan capacity effectively.

## Examples

### Basic info
Analyze the daily CPU utilization metrics of your OCI Core instances to understand usage patterns and performance trends. This can assist in optimizing resource allocation and identifying potential bottlenecks or underutilized instances.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_core_instance_metric_cpu_utilization_daily
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
  sample_count
from
  oci_core_instance_metric_cpu_utilization_daily
order by
  id,
  timestamp;
```

### CPU Over 80% average
Analyze the settings to understand instances where CPU utilization exceeds 80% on average. This is useful to identify potential performance issues and manage resource allocation effectively.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_core_instance_metric_cpu_utilization_daily
where average > 8
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  round(minimum,2) as min_cpu,
  round(maximum,2) as max_cpu,
  round(average,2) as avg_cpu,
  sample_count
from
  oci_core_instance_metric_cpu_utilization_daily
where average > 8
order by
  id,
  timestamp;
```