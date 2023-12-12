---
title: "Steampipe Table: oci_core_instance_metric_cpu_utilization_hourly - Query OCI Core Instance Metrics using SQL"
description: "Allows users to query CPU Utilization Metrics of OCI Core Instances on an hourly basis."
---

# Table: oci_core_instance_metric_cpu_utilization_hourly - Query OCI Core Instance Metrics using SQL

Oracle Cloud Infrastructure (OCI) Core instances provide secure, isolated compute environments for applications. They support a wide range of workloads and offer robust performance. CPU Utilization Metrics provide insights into the CPU usage of these instances, helping in performance monitoring and optimization.

## Table Usage Guide

The `oci_core_instance_metric_cpu_utilization_hourly` table provides insights into the CPU utilization of OCI Core instances on an hourly basis. As a system administrator or DevOps engineer, you can use this table to monitor CPU usage trends, identify potential performance bottlenecks, and make informed decisions about resource allocation and scaling. This table is particularly useful for maintaining optimal performance and ensuring efficient use of resources in your OCI environment.

## Examples

### Basic info
Explore the performance of your OCI core instances by analyzing their CPU utilization over time. This allows you to understand usage patterns and potentially optimize resource allocation.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_core_instance_metric_cpu_utilization_hourly
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
  oci_core_instance_metric_cpu_utilization_hourly
order by
  id,
  timestamp;
```

### CPU Over 80% average
Analyze the settings to understand instances where the CPU usage exceeds 80% on average. This can help in identifying potential performance issues and optimizing resource allocation.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_core_instance_metric_cpu_utilization_hourly
where average > 80
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
  oci_core_instance_metric_cpu_utilization_hourly
where average > 80
order by
  id,
  timestamp;
```