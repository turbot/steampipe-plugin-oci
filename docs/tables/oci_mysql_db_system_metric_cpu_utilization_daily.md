---
title: "Steampipe Table: oci_mysql_db_system_metric_cpu_utilization_daily - Query OCI MySQL DB System Metrics using SQL"
description: "Allows users to query daily CPU utilization metrics of MySQL DB Systems in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_mysql_db_system_metric_cpu_utilization_daily - Query OCI MySQL DB System Metrics using SQL

MySQL DB System is a fully managed database service in Oracle Cloud Infrastructure (OCI) that enables the deployment of cloud-native applications using the world's most popular open-source database. It provides a high-performance, seamless, and secure experience for developers and end users. This service is designed to handle any size workload from small to large scale applications.

## Table Usage Guide

The `oci_mysql_db_system_metric_cpu_utilization_daily` table provides insights into the daily CPU utilization metrics of MySQL DB Systems within Oracle Cloud Infrastructure (OCI). As a database administrator or a system engineer, explore CPU-specific details through this table, including average, maximum, and minimum CPU utilization. Utilize it to uncover information about CPU usage patterns, such as peak utilization times, and to aid in capacity planning and performance tuning.

## Examples

### Basic info
Explore the variations in CPU utilization of your MySQL database system over time. This analysis can help you understand system performance and plan for capacity needs.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization_daily
order by
  id,
  timestamp;
```

### CPU Over 80% average
Determine the instances where the average CPU utilization exceeds 80% to proactively manage system performance and avoid potential bottlenecks.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization_daily
where average > 80
order by
  id,
  timestamp;
```