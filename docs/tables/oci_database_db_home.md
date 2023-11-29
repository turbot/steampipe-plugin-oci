---
title: "Steampipe Table: oci_database_db_home - Query OCI Database DB Homes using SQL"
description: "Allows users to query OCI Database DB Homes."
---

# Table: oci_database_db_home - Query OCI Database DB Homes using SQL

A DB Home is a directory where Oracle Database software is installed. It can contain multiple databases and has specific directory structures and files. Within a DB System, you can have multiple DB Homes and each DB Home can run a different version of the Oracle Database software.

## Table Usage Guide

The `oci_database_db_home` table provides insights into DB Homes within Oracle Cloud Infrastructure Database service. As a Database administrator or a Cloud engineer, explore DB Home-specific details through this table, including its lifecycle state, compartment ID, and associated DB system ID. Utilize it to uncover information about DB Homes, such as their software version, time created, and the last patch history entry.

## Examples

### Basic info
Explore the lifecycle state and creation time of your Oracle Cloud Infrastructure databases to understand their current status and longevity. This can help in managing and planning resources effectively.

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_home;
```

### List db homes that are not available
Discover the database homes that are not currently available. This can be useful in identifying potential issues or disruptions in your database services.

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_home
where
  lifecycle_state <> 'AVAILABLE';
```