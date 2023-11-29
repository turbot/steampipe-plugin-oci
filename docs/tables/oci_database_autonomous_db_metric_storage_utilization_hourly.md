---
title: "Steampipe Table: oci_database_autonomous_db_metric_storage_utilization_hourly - Query OCI Database Autonomous Databases using SQL"
description: "Allows users to query hourly storage utilization metrics of Autonomous Databases in OCI."
---

# Table: oci_database_autonomous_db_metric_storage_utilization_hourly - Query OCI Database Autonomous Databases using SQL

Oracle Cloud Infrastructure's Autonomous Database is a fully managed, preconfigured database environment with two workload types available: Autonomous Transaction Processing and Autonomous Data Warehouse. The autonomous database is self-driving, self-securing, and self-repairing, which can automate all routine database tasks. It provides high performance, high availability, and strong security for any sized data warehouse workload.

## Table Usage Guide

The `oci_database_autonomous_db_metric_storage_utilization_hourly` table provides insights into the hourly storage utilization metrics of Autonomous Databases within Oracle Cloud Infrastructure's Database service. As a database administrator, explore database-specific details through this table, including storage consumed by the database, the timestamp of the metric, and the average storage used. Utilize it to monitor and manage your autonomous databases' storage utilization effectively, ensuring optimal performance and cost management.

## Examples

### Basic info
Explore the utilization of storage in your autonomous databases over time. This can help you monitor usage trends, identify periods of high demand, and plan for future capacity needs.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_database_autonomous_db_metric_storage_utilization_hourly
order by
  id,
  timestamp;
```

### Storage Utilization Over 80% average
Discover the segments that are utilizing more than 80% of their storage capacity on average. This information can help identify potential areas for optimization or expansion to ensure smooth operations.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_storage,
  round(maximum::numeric,2) as max_storage,
  round(average::numeric,2) as avg_storage,
  sample_count
from
  oci_database_autonomous_db_metric_storage_utilization_hourly
where average > 80
order by
  id,
  timestamp;
```