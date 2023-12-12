---
title: "Steampipe Table: oci_file_storage_mount_target - Query OCI File Storage Mount Targets using SQL"
description: "Allows users to query OCI File Storage Mount Targets."
---

# Table: oci_file_storage_mount_target - Query OCI File Storage Mount Targets using SQL

Oracle Cloud Infrastructure's File Storage service offers a durable, scalable, secure, enterprise-grade network file system. The service allows you to easily and securely share data across instances and containers. It provides a fully managed NFS file system that you can use with your applications without requiring you to install any file servers or hardware, or configure or manage NFS or SMB protocols.

## Table Usage Guide

The `oci_file_storage_mount_target` table provides insights into Mount Targets within Oracle Cloud Infrastructure's File Storage service. As a system administrator, explore Mount Target-specific details through this table, including export set details, lifecycle state, and associated metadata. Utilize it to uncover information about Mount Targets, such as those with specific export set details, the state of the Mount Target, and for monitoring and managing your file systems.

## Examples

### Basic info
Explore the status of your file storage mount targets to identify any that are not currently active. This can be useful in troubleshooting, ensuring optimal system performance, and managing resources effectively.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state as state,
  availability_domain,
  time_created
from
  oci_file_storage_mount_target;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state as state,
  availability_domain,
  time_created
from
  oci_file_storage_mount_target;
```

## List mount targets which are not active

```sql+postgres
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_mount_target
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_mount_target
where
  lifecycle_state <> 'ACTIVE';
```