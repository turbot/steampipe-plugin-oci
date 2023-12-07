---
title: "Steampipe Table: oci_database_autonomous_db_metric_cpu_utilization - Query OCI Database Autonomous DB Metric CPU Utilization using SQL"
description: "Allows users to query Autonomous DB Metric CPU Utilization in the OCI Database service."
---

# Table: oci_database_autonomous_db_metric_cpu_utilization - Query OCI Database Autonomous DB Metric CPU Utilization using SQL

An Autonomous DB Metric CPU Utilization is a resource in Oracle Cloud Infrastructure's (OCI) Database service. It provides metrics related to the CPU utilization of autonomous databases within an OCI environment. These metrics can be used to monitor the performance and efficiency of the databases, helping to identify any potential issues or areas for improvement.

## Table Usage Guide

The `oci_database_autonomous_db_metric_cpu_utilization` table provides insights into the CPU utilization of autonomous databases within the OCI Database service. As a Database Administrator, explore database-specific details through this table, including CPU utilization, average active sessions, and associated metadata. Utilize it to uncover information about autonomous databases, such as those with high CPU utilization, the performance of the databases, and the efficiency of resource usage.

## Examples

### Basic info
Explore the performance of your autonomous databases by analyzing their CPU utilization metrics. This query helps you assess the minimum, maximum, and average CPU usage over time, allowing you to optimize resources and improve database performance.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization
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
  oci_database_autonomous_db_metric_cpu_utilization
order by
  id,
  timestamp;
```

### CPU Over 80% average
Identify instances where the average CPU utilization exceeds 80% in your autonomous database. This allows you to monitor system performance and manage resources effectively.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization
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
  oci_database_autonomous_db_metric_cpu_utilization
where average > 80
order by
  id,
  timestamp;
```