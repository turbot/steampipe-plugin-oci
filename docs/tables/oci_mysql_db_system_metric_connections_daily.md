---
title: "Steampipe Table: oci_mysql_db_system_metric_connections_daily - Query OCI MySQL DB System Metrics using SQL"
description: "Allows users to query daily metrics of MySQL DB System Connections in OCI."
---

# Table: oci_mysql_db_system_metric_connections_daily - Query OCI MySQL DB System Metrics using SQL

MySQL DB System in Oracle Cloud Infrastructure (OCI) is a fully managed, highly available MySQL cloud service, enabling developers to rapidly develop and deploy secure, cloud native applications. It provides automated backups, failover, and repair to ensure high availability, and supports the latest MySQL features. It also offers advanced security features such as data encryption, network isolation, and MySQL Enterprise Firewall.

## Table Usage Guide

The `oci_mysql_db_system_metric_connections_daily` table provides insights into the daily metrics of MySQL DB System Connections in OCI. As a database administrator, you can use this table to monitor the connections to your MySQL databases on a daily basis, helping you to understand usage patterns and identify potential issues. This table can also be useful in capacity planning, by providing data on the number of connections over time.

## Examples

### Basic info
Explore the daily connection metrics of your MySQL database system in OCI to gain insights into its usage patterns and performance. This can help you identify any unusual activities or performance issues, and make necessary adjustments to optimize the system's efficiency.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_connections_daily
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
  oci_mysql_db_system_metric_connections_daily
order by
  id,
  timestamp;
```

### Active connection statistics
Explore the fluctuation of active connections to your MySQL database system over time. This information can help you understand your system's usage patterns and manage resource allocation more effectively.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_daily
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
  oci_mysql_db_system_metric_connections_daily
where metric_name = 'ActiveConnections'
order by
  id,
  timestamp;
```

### Current connection statistics
Analyze the daily connection statistics of your MySQL database system to gain insights into its usage patterns and performance. This analysis could help in optimizing the system configuration for improved performance and resource management.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_conn,
  round(maximum::numeric,2) as max_conn,
  round(average::numeric,2) as avg_conn,
  sample_count
from
  oci_mysql_db_system_metric_connections_daily
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
  oci_mysql_db_system_metric_connections_daily
where metric_name = 'CurrentConnections'
order by
  id,
  timestamp;
```