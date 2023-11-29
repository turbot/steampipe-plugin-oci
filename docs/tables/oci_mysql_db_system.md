---
title: "Steampipe Table: oci_mysql_db_system - Query OCI MySQL Database Systems using SQL"
description: "Allows users to query detailed information about MySQL Database Systems within Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_mysql_db_system - Query OCI MySQL Database Systems using SQL

A MySQL Database System is a managed service provided by Oracle Cloud Infrastructure (OCI) that makes it easy to set up, operate, and scale MySQL deployments in the cloud. It offers cost-efficient and resizable capacity while automating time-consuming administration tasks such as hardware provisioning, database setup, patching and backups. It frees you to focus on your applications so you can give them the fast performance, high availability, security and compatibility they need.

## Table Usage Guide

The `oci_mysql_db_system` table provides insights into MySQL Database Systems within Oracle Cloud Infrastructure (OCI). As a Database Administrator, you can explore system-specific details through this table, including configurations, network settings, and associated metadata. Utilize it to uncover information about systems, such as their current lifecycle state, the associated configurations, and the details of the subnet in which the system resides.

## Examples

### Basic info
Explore the lifecycle status and creation date of your MySQL database systems to understand their operational state and longevity. This can be useful in managing and assessing the overall health of your databases.

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system;
```

### List DB systems that are not active
Discover the segments that consist of DB systems that are not currently active. This is useful in identifying and managing inactive resources within your MySQL database.

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system
where
  lifecycle_state <> 'ACTIVE';
```

### List DB systems with backups not enabled
Discover the segments that have active database systems without backup enabled, enabling you to identify potential vulnerabilities and risks in your data management. This is crucial for maintaining data integrity and implementing disaster recovery plans.

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system
where
  lifecycle_state = 'ACTIVE'
  and backup_policy -> 'isEnabled' <> 'true';
```

### List the CPU and RAM configuration of DB systems
Explore the active DB systems to analyze their processing and memory capabilities. This helps in understanding the hardware resource allocation, supporting efficient system management and performance optimization.

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  cpu_core_count,
  memory_size_in_gbs
from
  oci_mysql_db_system
where
  lifecycle_state = 'ACTIVE';
```