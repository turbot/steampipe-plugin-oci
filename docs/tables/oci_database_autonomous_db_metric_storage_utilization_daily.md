---
title: "Steampipe Table: oci_database_autonomous_db_metric_storage_utilization_daily - Query OCI Database Autonomous Databases using SQL"
description: "Allows users to query daily metrics related to storage utilization for OCI Autonomous Databases."
---

# Table: oci_database_autonomous_db_metric_storage_utilization_daily - Query OCI Database Autonomous Databases using SQL

Oracle Cloud Infrastructure's Autonomous Database is a fully managed, preconfigured database environment with two workload types available - Autonomous Transaction Processing and Autonomous Data Warehouse. The autonomous database is self-driving, self-securing, and self-repairing, which can automate all routine database tasks. It uses machine learning algorithms to enable unprecedented availability, high performance, and security, delivering a significantly lower cost of ownership.

## Table Usage Guide

The `oci_database_autonomous_db_metric_storage_utilization_daily` table provides insights into the daily storage utilization metrics of OCI Autonomous Databases. As a database administrator, you can use this table to monitor and analyze storage usage trends and patterns over time. This can help in proactive capacity planning, ensuring optimal performance, and avoiding potential storage-related issues.

## Examples

### Basic info
Explore the patterns of storage utilization in your Oracle Cloud Infrastructure's autonomous databases. This query helps you assess the minimum, maximum, and average storage used daily, aiding in efficient resource management.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_db_metric_storage_utilization_daily
order by
  id,
  timestamp;
```

### Storage Utilization Over 80% average
Determine the areas in which storage utilization exceeds an average of 80%. This query is useful for identifying potential issues with storage management and can aid in optimizing resource allocation.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_storage,
  round(maximum::numeric,2) as max_storage,
  round(average::numeric,2) as avg_storage,
  sample_count
from
  oci_database_autonomous_db_metric_storage_utilization_daily
where average > 80
order by
  id,
  timestamp;
```