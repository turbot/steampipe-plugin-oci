---
title: "Steampipe Table: oci_database_autonomous_db_metric_cpu_utilization_daily - Query OCI Database Autonomous Databases using SQL"
description: "Allows users to query daily CPU Utilization Metrics of Autonomous Databases in Oracle Cloud Infrastructure."
---

# Table: oci_database_autonomous_db_metric_cpu_utilization_daily - Query OCI Database Autonomous Databases using SQL

Oracle Autonomous Database is a fully managed, preconfigured database environment with two workload types available, Autonomous Transaction Processing and Autonomous Data Warehouse. The autonomous database is built on Oracle Database 19c and includes features like automatic indexing, AI-driven tuning, and automated data movement. It is designed to support all standard SQL and business intelligence (BI) tools and provides all of the performance of the market-leading Oracle Database in an environment that is tuned and optimized for data warehouse workloads.

## Table Usage Guide

The `oci_database_autonomous_db_metric_cpu_utilization_daily` table provides insights into daily CPU Utilization Metrics of Autonomous Databases within Oracle Cloud Infrastructure. As a Database Administrator, explore database-specific details through this table, including CPU utilization, average active sessions, and associated metadata. Utilize it to monitor and optimize the performance of your Autonomous Databases, such as identifying databases with high CPU utilization, and making informed decisions on resource allocation.

## Examples

### Basic info
Analyze the daily CPU utilization metrics of your autonomous databases to understand their performance trends and resource consumption. This can assist in optimizing resource allocation, identifying peak usage times, and planning capacity for better database management.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization_daily
order by
  id,
  timestamp;
```

### CPU Over 80% average
Identify instances where the average CPU utilization exceeds 80% in a day to monitor system performance and prevent potential overloads or slowdowns. This allows for proactive resource management and can help maintain optimal system operation.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_database_autonomous_db_metric_cpu_utilization_daily
where average > 80
order by
  id,
  timestamp;
```