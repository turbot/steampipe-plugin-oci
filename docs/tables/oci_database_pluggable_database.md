---
title: "Steampipe Table: oci_database_pluggable_database - Query OCI Database Pluggable Databases using SQL"
description: "Allows users to query OCI Database Pluggable Databases."
---

# Table: oci_database_pluggable_database - Query OCI Database Pluggable Databases using SQL

A Pluggable Database (PDB) is a portable collection of schemas, schema objects, and non-schema objects that appears to an Oracle Net client as a non-CDB. All Oracle databases before Oracle Database 12c were non-CDBs. The data and code of a PDB are physically stored in common tablespaces, which are owned by the CDB root.

## Table Usage Guide

The `oci_database_pluggable_database` table provides insights into Pluggable Databases within Oracle Cloud Infrastructure Database service. As a database administrator, you can explore specific details about each PDB through this table, including its configuration, status, and associated metadata. Utilize it to uncover information about PDBs, such as their storage size, the number of users, and the version of the Oracle Database software.

## Examples

### Basic info
Explore which databases are in different stages of their lifecycle and when they were created. This allows you to assess the overall health and status of your databases.

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database;
```

### List failed pluggable databases
This query can be used to identify any pluggable databases in the Oracle Cloud Infrastructure that have failed to initialize. By doing so, it allows users to quickly pinpoint and address any issues that may be disrupting their database operations.

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database
where
  lifecycle_state = 'FAILED';
```

### List pluggable databases older than 90 days
Discover the segments that consist of pluggable databases older than 90 days. This is useful for identifying potential areas for system optimization or data archiving.

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

### List unrestricted pluggable databases
Explore which pluggable databases are unrestricted in your Oracle Cloud Infrastructure. This can be useful to identify potential security vulnerabilities or for general database management purposes.

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  is_restricted
from
  oci_database_pluggable_database
where
  not is_restricted;
```