---
title: "Steampipe Table: oci_database_autonomous_db_metric_storage_utilization - Query OCI Database Autonomous Databases using SQL"
description: "Allows users to query metrics related to storage utilization for Oracle Cloud Infrastructure (OCI) Database Autonomous Databases."
---

# Table: oci_database_autonomous_db_metric_storage_utilization - Query OCI Database Autonomous Databases using SQL

Oracle Cloud Infrastructure (OCI) Database Autonomous Database is a fully managed, pre-configured database environment with two workload types available - Autonomous Transaction Processing and Autonomous Data Warehouse. The environment uses machine learning algorithms to enable automation of database tuning, security, backups, updates, and other routine management tasks traditionally performed by database administrators. It offers high-performance, reliability, and seamless scalability with a broad suite of developer tools.

## Table Usage Guide

The `oci_database_autonomous_db_metric_storage_utilization` table provides insights into the storage utilization metrics of OCI Database Autonomous Databases. As a database administrator or data analyst, you can use this table to monitor and manage storage utilization, enabling you to optimize database performance and resource allocation. It can also be useful for auditing and compliance purposes, helping you ensure that storage usage aligns with organizational policies and industry regulations.

## Examples

### Basic info
Explore the storage utilization of your autonomous database to gain insights into usage trends over time. This can help in planning storage capacity and maintaining optimal performance.

```sql+postgres
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_database_metric_storage_utilization
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
  oci_database_autonomous_database_metric_storage_utilization
order by
  id,
  timestamp;
```

### Storage Utilization Over 80% average
Identify instances where the average storage utilization of your autonomous database in OCI exceeds 80%. This can help in managing resources effectively by pinpointing databases that might need storage optimization or capacity upgrades.

```sql+postgres
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_storage,
  round(maximum::numeric,2) as max_storage,
  round(average::numeric,2) as avg_storage,
  sample_count
from
  oci_database_autonomous_database_metric_storage_utilization
where average > 80
order by
  id,
  timestamp;
```

```sql+sqlite
select
  id,
  timestamp,
  round(minimum,2) as min_storage,
  round(maximum,2) as max_storage,
  round(average,2) as avg_storage,
  sample_count
from
  oci_database_autonomous_database_metric_storage_utilization
where average > 80
order by
  id,
  timestamp;
```