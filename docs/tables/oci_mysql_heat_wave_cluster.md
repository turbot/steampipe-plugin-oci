---
title: "Steampipe Table: oci_mysql_heat_wave_cluster - Query OCI MySQL HeatWave Clusters using SQL"
description: "Allows users to query data from OCI MySQL HeatWave Clusters."
---

# Table: oci_mysql_heat_wave_cluster - Query OCI MySQL HeatWave Clusters using SQL

MySQL HeatWave is an integrated, high-performance analytic engine for the MySQL Database Service in Oracle Cloud Infrastructure (OCI). It accelerates MySQL performance by orders of magnitude for analytics and scales to millions of queries per second. MySQL HeatWave provides a single, unified platform for both transaction processing and analytics, eliminating the need for ETL, data duplication, and data movement.

## Table Usage Guide

The `oci_mysql_heat_wave_cluster` table provides insights into MySQL HeatWave Clusters within Oracle Cloud Infrastructure. As a database administrator or data analyst, explore cluster-specific details through this table, including cluster status, capacity, and associated metadata. Utilize it to uncover information about clusters, such as their configuration settings, performance metrics, and the verification of operational status.

## Examples

### Basic info
Explore the lifecycle state and creation time of your MySQL HeatWave clusters in Oracle Cloud Infrastructure. This can help you understand their status and longevity, which is useful for tracking usage and managing resources.

```sql+postgres
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster;
```

```sql+sqlite
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster;
```

### List failed heat wave clusters
Explore which MySQL Heat Wave clusters have failed to understand potential issues in your database system. This can help in troubleshooting and improving the system's reliability.

```sql+postgres
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster
where
  lifecycle_state = 'FAILED';
```

```sql+sqlite
select
  db_system_id,
  lifecycle_state as state,
  time_created
from
  oci_mysql_heat_wave_cluster
where
  lifecycle_state = 'FAILED';
```

### List heat wave clusters older than 90 days
Identify older heat wave clusters in your MySQL system to assess their lifecycle state and size. This is particularly useful for managing system resources and maintaining optimal database performance.

```sql+postgres
select
  db_system_id,
  lifecycle_state as state,
  time_created,
  cluster_size
from
  oci_mysql_heat_wave_cluster
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

```sql+sqlite
select
  db_system_id,
  lifecycle_state as state,
  time_created,
  cluster_size
from
  oci_mysql_heat_wave_cluster
where
  time_created <= date('now','-90 day')
order by
  time_created;
```