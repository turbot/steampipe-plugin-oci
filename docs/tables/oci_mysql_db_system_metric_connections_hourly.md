---
title: "Steampipe Table: oci_mysql_db_system_metric_connections_hourly - Query OCI MySQL DB Systems using SQL"
description: "Allows users to query hourly metrics for connections to MySQL DB Systems in OCI."
---

# Table: oci_mysql_db_system_metric_connections_hourly - Query OCI MySQL DB Systems using SQL

Oracle Cloud Infrastructure's MySQL DB System is a fully managed database service that allows you to deploy cloud-native applications using the world's most popular open source database. It is built on MySQL Enterprise Edition and powered by the Oracle Cloud, providing a simple, automated, integrated and enterprise-ready MySQL cloud service. The service enables organizations to increase business agility and reduce costs.

## Table Usage Guide

The `oci_mysql_db_system_metric_connections_hourly` table provides insights into the hourly connection metrics of MySQL DB Systems within Oracle Cloud Infrastructure (OCI). As a database administrator, you can explore connection-specific details through this table, including the number of successful, rejected, and total connection attempts. Utilize it to monitor connection trends, identify potential issues, and ensure optimal performance of your MySQL DB Systems.

## Examples

### Basic info
Explore the hourly metrics of MySQL database system connections to understand usage patterns and performance. This can help in optimizing resource allocation and identifying potential issues in real-time.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
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
  oci_mysql_db_system_metric_connections_hourly
order by
  id,
  timestamp;
```

### Active connection statistics
Analyze the active connection statistics in your MySQL database system to gain insights into the minimum, maximum, and average connections over time. This information can help you understand usage patterns and optimize resource allocation for improved performance.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
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
  oci_mysql_db_system_metric_connections_hourly
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

### Current connection statistics
Explore the varying connection statistics in your MySQL database system over time to better understand usage patterns and anticipate potential issues. This query helps in monitoring the minimum, maximum, and average number of connections to optimize performance and resource allocation.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_hourly
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
  oci_mysql_db_system_metric_connections_hourly
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```