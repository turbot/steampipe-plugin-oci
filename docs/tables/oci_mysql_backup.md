---
title: "Steampipe Table: oci_mysql_backup - Query OCI MySQL Backups using SQL"
description: "Allows users to query MySQL Backups in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_mysql_backup - Query OCI MySQL Backups using SQL

MySQL Backups in Oracle Cloud Infrastructure (OCI) are automated or manual backups of your MySQL DB System. These backups are stored in OCI Object Storage and are used to restore a DB System to a specific point in time. The backup data includes the data volume and the system volume.

## Table Usage Guide

The `oci_mysql_backup` table provides insights into MySQL backups within Oracle Cloud Infrastructure (OCI). As a DBA or DevOps engineer, explore backup-specific details through this table, including backup type, lifecycle state, and associated metadata. Utilize it to uncover information about backups, such as those with specific creation and expiry times, the DB system associated with the backup, and the verification of backup configurations.

## Examples

### Basic info
Analyze your MySQL backups in Oracle Cloud Infrastructure to understand their current lifecycle state, backup type, and MySQL version. This can help manage backups more efficiently and ensure they're configured correctly for your needs.

```sql+postgres
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  backup_type,
  mysql_version
from
  oci_mysql_backup;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  backup_type,
  mysql_version
from
  oci_mysql_backup;
```

### List manual backups
Explore which MySQL backups have been manually created. This is useful for understanding how many backups have been created by user intervention, which can help in managing resources and planning future automated backups.

```sql+postgres
select
  display_name,
  id,
  creation_type
from
  oci_mysql_backup
where
  creation_type = 'MANUAL';
```

```sql+sqlite
select
  display_name,
  id,
  creation_type
from
  oci_mysql_backup
where
  creation_type = 'MANUAL';
```

### List backups with retention days less than 90 days
Determine areas in which backups have a retention period of less than 90 days. This can be useful for identifying potential vulnerabilities or non-compliance with data retention policies.

```sql+postgres
select
  display_name,
  id,
  retention_in_days
from
  oci_mysql_backup
where
  retention_in_days < 90;
```

```sql+sqlite
select
  display_name,
  id,
  retention_in_days
from
  oci_mysql_backup
where
  retention_in_days < 90;
```

### Count of backups per DB system
Explore the frequency of backups for each database system to understand how often data protection measures are being implemented. This can help in assessing the robustness of your data recovery strategy.

```sql+postgres
select
  db_system_id,
  count(*) as backup_count
from
  oci_mysql_backup
group by
  db_system_id;
```

```sql+sqlite
select
  db_system_id,
  count(*) as backup_count
from
  oci_mysql_backup
group by
  db_system_id;
```