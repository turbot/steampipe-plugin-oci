---
title: "Steampipe Table: oci_core_instance_metric_cpu_utilization - Query OCI Core Instance Metrics using SQL"
description: "Allows users to query CPU Utilization Metrics for OCI Core Instances."
---

# Table: oci_core_instance_metric_cpu_utilization - Query OCI Core Instance Metrics using SQL

OCI Core Instance Metrics are a part of Oracle Cloud Infrastructure's Monitoring service. These metrics provide real-time data about the performance of your instances. CPU Utilization Metrics specifically provide data about the percentage of total CPU resources that an instance uses.

## Table Usage Guide

The `oci_core_instance_metric_cpu_utilization` table provides insights into CPU Utilization Metrics for OCI Core Instances. As a system administrator or DevOps engineer, you can use this table to monitor and manage the performance of your instances. This table can be particularly useful for identifying instances that are under heavy load or are not utilizing their CPU resources efficiently.

## Examples

### Basic info
Explore the performance of your virtual machines by analyzing their CPU utilization metrics. This allows you to assess the efficiency of your resources and make informed decisions about scaling or resource allocation.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_core_instance_metric_cpu_utilization
order by
  id,
  timestamp;
```

### CPU Over 80% average
Determine the instances where the CPU utilization exceeds 80% on average. This is useful for identifying potential performance issues and ensuring optimal resource management.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_core_instance_metric_cpu_utilization
where average > 80
order by
  id,
  timestamp;
```