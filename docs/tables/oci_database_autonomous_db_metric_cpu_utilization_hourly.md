---
title: "Steampipe Table: oci_database_autonomous_db_metric_cpu_utilization_hourly - Query OCI Database Autonomous Databases CPU Utilization Metrics using SQL"
description: "Allows users to query Autonomous Databases CPU Utilization Metrics in Oracle Cloud Infrastructure (OCI) Database service."
---

# Table: oci_database_autonomous_db_metric_cpu_utilization_hourly - Query OCI Database Autonomous Databases CPU Utilization Metrics using SQL

Autonomous Databases in Oracle Cloud Infrastructure Database service are self-driving, self-securing, and self-repairing databases that automate key management processes, including patching, tuning, backups, and upgrades. These databases provide CPU Utilization Metrics, which provide insights into the CPU usage of the databases. These metrics are important for monitoring the performance and ensuring the smooth operation of the databases.

## Table Usage Guide

The `oci_database_autonomous_db_metric_cpu_utilization_hourly` table provides insights into the CPU utilization metrics of Autonomous Databases in OCI Database service. As a Database Administrator or DevOps engineer, you can leverage this table to monitor and manage the performance of your Autonomous Databases, including identifying high CPU usage periods, planning for capacity, and optimizing resource allocation. It serves as a valuable tool for maintaining the efficiency and reliability of your databases.

## Examples

### Basic info
Explore the patterns of CPU utilization over time in your Oracle Cloud Infrastructure database. This can help you understand resource usage trends and plan for capacity adjustments accordingly.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization_hourly
order by
  id,
  timestamp;
```

### CPU Over 80% average
Explore instances where the average CPU usage exceeds 80% to identify potential performance issues and optimize resource allocation. This can be particularly beneficial in managing resources efficiently and preventing system slowdowns or crashes due to high CPU utilization.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization_hourly
where average > 80
order by
  id,
  timestamp;
```