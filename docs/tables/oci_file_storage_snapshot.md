---
title: "Steampipe Table: oci_file_storage_snapshot - Query OCI File Storage Snapshots using SQL"
description: "Allows users to query OCI File Storage Snapshots."
---

# Table: oci_file_storage_snapshot - Query OCI File Storage Snapshots using SQL

Oracle Cloud Infrastructure's File Storage service provides a durable, scalable, secure, enterprise-grade network file system. It allows you to dynamically grow and shrink your file system without the need to provision capacity ahead of time. File Storage supports the Network File System (NFS) protocol.

## Table Usage Guide

The `oci_file_storage_snapshot` table provides insights into snapshots within OCI File Storage. As a DevOps engineer, explore snapshot-specific details through this table, including the snapshot's lifecycle state, time created, and associated metadata. Utilize it to uncover information about snapshots, such as those with a specific lifecycle state or created at a specific time.

## Examples

### Basic info
1. "Explore the basic information about your file storage snapshots, such as their names, IDs, states, creation times, provenance IDs, and regions."
2. "Determine the number of snapshots created per file system to understand the distribution and usage of your storage resources."
3. "Identify the snapshots that have been used as clone sources, helping you track the propagation of your data across your storage environment.

```sql+postgres
select
  name,
  id,
  lifecycle_state as state,
  time_created,
  provenance_id,
  region
from
  oci_file_storage_snapshot;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state as state,
  time_created,
  provenance_id,
  region
from
  oci_file_storage_snapshot;
```


## Count of snapshots created per file system

```sql+postgres
select
  file_system_id,
  count(*) as snapshots_count
from
  oci_file_storage_snapshot
group by
  file_system_id;
```

```sql+sqlite
select
  file_system_id,
  count(*) as snapshots_count
from
  oci_file_storage_snapshot
group by
  file_system_id;
```


## List cloned snapshots

```sql+postgres
select
  name,
  id,
  is_clone_source
from
  oci_file_storage_snapshot
where
  is_clone_source;
```

```sql+sqlite
select
  name,
  id,
  is_clone_source
from
  oci_file_storage_snapshot
where
  is_clone_source = 1;
```