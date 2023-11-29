---
title: "Steampipe Table: oci_mysql_db_system_metric_cpu_utilization_hourly - Query OCI MySQL DB System Metrics using SQL"
description: "Allows users to query MySQL DB System CPU Utilization Metrics on an hourly basis."
---

# Table: oci_mysql_db_system_metric_cpu_utilization_hourly - Query OCI MySQL DB System Metrics using SQL

Oracle Cloud Infrastructure's MySQL DB System is a fully managed, scalable MySQL relational database service that enables organizations to deploy cloud-native applications. It offers a secure, automated, and extensible platform for running MySQL applications. This service provides the ability to monitor CPU utilization metrics on an hourly basis.

## Table Usage Guide

The `oci_mysql_db_system_metric_cpu_utilization_hourly` table provides insights into the CPU utilization metrics of MySQL DB Systems within Oracle Cloud Infrastructure (OCI). As a database administrator, you can explore these metrics to understand the CPU usage patterns and performance of your MySQL DB Systems. Utilize it to uncover information about CPU usage trends, identify peak usage times, and plan capacity accordingly.

## Examples

### Basic info
Analyze the CPU utilization of your MySQL database system over the past hour. This allows you to identify periods of high demand and optimize system performance accordingly.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization_hourly
order by
  id,
  timestamp;
```

### CPU Over 80% average
Analyze the settings to understand instances where CPU utilization is consistently high, exceeding 80% on average. This is useful for identifying potential performance bottlenecks and planning for capacity upgrades.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization_hourly
where average > 80
order by
  id,
  timestamp;
```