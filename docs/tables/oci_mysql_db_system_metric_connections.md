---
title: "Steampipe Table: oci_mysql_db_system_metric_connections - Query OCI MySQL DB System Metrics using SQL"
description: "Allows users to query MySQL DB System Metrics related to connections."
---

# Table: oci_mysql_db_system_metric_connections - Query OCI MySQL DB System Metrics using SQL

Oracle Cloud Infrastructure (OCI) MySQL Database Service is a fully managed relational database service that enables organizations to deploy cloud-native applications using the world’s most popular open-source database. It is developed, managed, and supported by the MySQL team at Oracle. The service is built on the highly available and secure OCI platform and is fully compatible with on-premises MySQL databases.

## Table Usage Guide

The `oci_mysql_db_system_metric_connections` table provides insights into connection metrics within OCI MySQL Database Service. As a database administrator, explore connection-specific details through this table, including total connections, successful connections, and rejected connections. Utilize it to uncover information about database connections, such as those with high frequency, the successful and failed connection attempts, and the verification of connection policies.

## Examples

### Basic info
Explore the performance metrics of your MySQL database system over time. This query helps you assess the minimum, maximum, and average number of connections, along with the total sample count, allowing you to better understand your system's utilization and manage its capacity effectively.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_connections
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
  oci_mysql_db_system_metric_connections
order by
  id,
  timestamp;
```

### Active connection statistics
Analyze the statistics of active connections to gain insights into the performance and usage patterns of your MySQL database system. This can help you understand peak usage times, manage resources more effectively, and troubleshoot potential issues.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  round(minimum,2) as min_conn,
  round(maximum,2) as max_conn,
  round(average,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

### Current connection statistics
Analyze the current connection statistics from your MySQL database system in Oracle Cloud Infrastructure (OCI) to understand the minimum, maximum, and average connections over a specified period. This can help you manage your database connections more effectively, ensuring optimal performance and resource utilization.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  round(minimum,2) as min_conn,
  round(maximum,2) as max_conn,
  round(average,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```