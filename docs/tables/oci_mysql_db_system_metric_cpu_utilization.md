---
title: "Steampipe Table: oci_mysql_db_system_metric_cpu_utilization - Query OCI MySQL Database Service DB System Metrics using SQL"
description: "Allows users to query DB System Metrics for CPU Utilization in the OCI MySQL Database Service."
---

# Table: oci_mysql_db_system_metric_cpu_utilization - Query OCI MySQL Database Service DB System Metrics using SQL

The OCI MySQL Database Service is a fully managed relational database service that enables organizations to deploy cloud-native applications using the world's most popular open-source database. It is built on MySQL Enterprise Edition and powered by Oracle Cloud, providing a cost-effective and high-performance solution for deploying applications in the cloud. DB System Metrics for CPU Utilization provides valuable insights into the CPU usage of your MySQL Database System.

## Table Usage Guide

The `oci_mysql_db_system_metric_cpu_utilization` table provides insights into the CPU utilization of DB Systems in the OCI MySQL Database Service. As a database administrator, you can explore detailed metrics about CPU usage through this table, including the percentage of CPU utilization, the timestamp of the data, and the average, maximum, and minimum values over a specified time period. Utilize it to monitor and optimize the performance of your MySQL Database Systems in the cloud.

## Examples

### Basic info
Explore the performance metrics of a MySQL database system to understand its CPU utilization. This allows you to assess the system's efficiency and optimize its performance by identifying the instances of minimum, maximum, and average CPU usage over a period of time.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization
order by
  id,
  timestamp;
```

### CPU Over 80% average
Identify instances where the average CPU usage surpasses 80%. This query is useful in monitoring system health and preventing potential overloads or crashes due to high CPU usage.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_cpu,
  round(maximum::numeric,2) as max_cpu,
  round(average::numeric,2) as avg_cpu,
  sample_count
from
  oci_mysql_db_system_metric_cpu_utilization
where average > 80
order by
  id,
  timestamp;
```