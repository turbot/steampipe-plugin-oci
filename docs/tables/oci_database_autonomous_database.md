---
title: "Steampipe Table: oci_database_autonomous_database - Query OCI Database Autonomous Databases using SQL"
description: "Allows users to query Autonomous Databases within the Oracle Cloud Infrastructure Database service."
---

# Table: oci_database_autonomous_database - Query OCI Database Autonomous Databases using SQL

The Autonomous Database is a feature of the Oracle Cloud Infrastructure Database service. It automates management tasks such as tuning, patching, backups and more, allowing developers to focus on higher-value tasks. The Autonomous Database supports a variety of workloads, including OLTP, data warehousing, and JSON document store.

## Table Usage Guide

The `oci_database_autonomous_database` table provides insights into Autonomous Databases within Oracle Cloud Infrastructure Database service. As a database administrator or developer, explore specific details about these databases, including their configurations, statuses, and associated metadata. Utilize this table to manage and monitor your databases, ensuring they are optimized, secure, and functioning as expected.

## Examples

### Basic info
Gain insights into the lifecycle state and creation time of your autonomous databases to better understand their status and duration of existence. This is particularly useful for database management and auditing purposes.

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database;
```

### List databases that are not available
Discover the databases that are currently not available. This can be useful to identify potential issues or disruptions in your database services.

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database
where
  lifecycle_state <> 'AVAILABLE';
```

### List databases with a data storage size greater than 1024 GB
Discover the databases that have a storage size exceeding 1024 GB. This query is useful to monitor and manage your database storage, aiding in efficient resource allocation and preventing potential storage shortages.

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database
where
  data_storage_size_in_gbs > 1024;
```

### Get KMS key details for the databases
Discover the encryption key details for your databases to gain insights into their security measures. This can be particularly useful for assessing the strength of your data protection and identifying areas for potential improvement.

```sql
select
  d.db_name,
  d.display_name,
  d.kms_key_id,
  k.name as key_name,
  k.algorithm as key_algorithm,
  k.current_key_version,
  k.protection_mode
from
  oci_database_autonomous_database as d,
  oci_kms_key as k
where
  k.id = d.kms_key_id;
```

### Get KMS vault details for the databases
Determine the areas in which your databases are connected to specific KMS vaults. This can help you understand the security measures in place for your data, as well as identify potential areas for improvement or optimization.

```sql
select
  d.db_name,
  d.display_name,
  d.vault_id,
  v.display_name as vault_display_name,
  v.crypto_endpoint,
  v.vault_type,
  v.management_endpoint
from
  oci_database_autonomous_database as d,
  oci_kms_vault as v
where
  v.id = d.vault_id;
```