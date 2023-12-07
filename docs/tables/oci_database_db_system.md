---
title: "Steampipe Table: oci_database_db_system - Query OCI Database Service DB Systems using SQL"
description: "Allows users to query DB Systems in the OCI Database Service."
---

# Table: oci_database_db_system - Query OCI Database Service DB Systems using SQL

DB Systems in OCI Database Service provide a fully managed database service that lets you create, manage, and scale Oracle Database instances in the cloud. It is designed to support any size workload from small dev/test to large production deployments. The service provides automated backup, recovery, and patching while providing high availability and security.

## Table Usage Guide

The `oci_database_db_system` table provides insights into DB Systems within the Oracle Cloud Infrastructure Database Service. As a database administrator or DevOps engineer, you can explore DB System-specific details through this table, including its configuration, status, and associated resources. Utilize this table to manage and monitor your OCI Database Service resources effectively, such as identifying DB Systems that require patching or understanding the distribution of DB Systems across different compartments.

## Examples

### Basic info
Explore which lifecycle states your Oracle Cloud Infrastructure databases are in and when they were created. This can help you manage your databases by identifying those that are active, inactive, or in transition.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system;
```

### List db systems that are not available
Explore which database systems are not currently available. This can be particularly useful in identifying potential issues or disruptions in your database environment.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  lifecycle_state <> 'AVAILABLE';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  lifecycle_state <> 'AVAILABLE';
```

### List db systems with data storage size greater than 1024 GB
Explore which database systems have a data storage size exceeding 1024 GB to determine if there's a need for storage optimization or capacity planning. This can be beneficial in managing resources and preventing potential system overloads.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  data_storage_size_in_gbs > 1024;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  data_storage_size_in_gbs > 1024;
```