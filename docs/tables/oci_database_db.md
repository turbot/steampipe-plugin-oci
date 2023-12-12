---
title: "Steampipe Table: oci_database_db - Query OCI Database DBs using SQL"
description: "Allows users to query OCI Database DBs."
---

# Table: oci_database_db - Query OCI Database DBs using SQL

Oracle Cloud Infrastructure (OCI) Database service provides a fully managed, scalable, and highly available database solution. It supports a variety of database workloads, including Online Transaction Processing (OLTP), data warehousing (DW), and mixed workloads (OLTP and DW). The service is built on Oracle Exadata Database Machine, ensuring high performance and availability.

## Table Usage Guide

The `oci_database_db` table provides insights into databases within Oracle Cloud Infrastructure (OCI) Database service. As a database administrator, explore database-specific details through this table, including the lifecycle state, database name, and associated metadata. Utilize it to uncover information about databases, such as those in a terminated state, the relationships between databases, and the verification of database configurations.

## Examples

### Basic info
Explore which databases are active or inactive by assessing their lifecycle state and when they were created. This information can be useful in managing and maintaining your database resources effectively.

```sql+postgres
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db;
```

```sql+sqlite
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db;
```

### List databases that are not available
Discover the databases that are not currently available. This is useful in identifying potential issues or maintenance requirements within your database system.

```sql+postgres
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db
where
  lifecycle_state <> 'AVAILABLE';
```

```sql+sqlite
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db
where
  lifecycle_state <> 'AVAILABLE';
```