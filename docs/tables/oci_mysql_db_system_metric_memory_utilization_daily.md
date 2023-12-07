---
title: "Steampipe Table: oci_mysql_db_system_metric_memory_utilization_daily - Query OCI MySQL Database Services DB System Metrics using SQL"
description: "Allows users to query daily memory utilization metrics for MySQL DB Systems in Oracle Cloud Infrastructure."
---

# Table: oci_mysql_db_system_metric_memory_utilization_daily - Query OCI MySQL Database Services DB System Metrics using SQL

The Oracle Cloud Infrastructure (OCI) MySQL Database Service is a fully managed database service that lets developers and database administrators focus on application development and business needs rather than infrastructure operations. This service provides the ability to deploy MySQL Server-based applications in the cloud, on a platform that's built for scale, performance, and security. The MySQL Database Service is the only public cloud service built on MySQL Enterprise Edition.

## Table Usage Guide

The `oci_mysql_db_system_metric_memory_utilization_daily` table provides insights into the daily memory utilization metrics of MySQL DB Systems within Oracle Cloud Infrastructure (OCI). As a Database Administrator or Developer, you can explore detailed memory usage statistics through this table, including total memory, used memory, and free memory. Utilize it to monitor and optimize your MySQL DB Systems' performance and resource usage, ensuring optimal operation and cost-effectiveness.

## Examples

### Basic info
Explore the daily memory utilization patterns of your MySQL database system. This helps in understanding memory usage trends and can aid in capacity planning and performance optimization.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization_daily
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
  oci_mysql_db_system_metric_memory_utilization_daily
order by
  id,
  timestamp;
```

### Memory Utilization Over 80% average
This query helps to pinpoint instances where memory utilization has exceeded 80% on average. This is useful for identifying potential system performance issues and planning for capacity upgrades.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_memory,
  round(maximum::numeric,2) as max_memory,
  round(average::numeric,2) as avg_memory,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization_daily
where average > 80
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  round(minimum,2) as min_memory,
  round(maximum,2) as max_memory,
  round(average,2) as avg_memory,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization_daily
where average > 80
order by
  id,
  timestamp;
```