---
title: "Steampipe Table: oci_file_storage_file_system - Query OCI File Storage Systems using SQL"
description: "Allows users to query OCI File Storage Systems."
---

# Table: oci_file_storage_file_system - Query OCI File Storage Systems using SQL

Oracle Cloud Infrastructure (OCI) File Storage service provides a durable, scalable, secure, enterprise-grade network file system. You can connect to a File Storage service file system from any bare metal, virtual machine, or container instance in your Virtual Cloud Network (VCN). The File Storage service supports the Network File System (NFS) protocol.

## Table Usage Guide

The `oci_file_storage_file_system` table provides insights into file systems within Oracle Cloud Infrastructure (OCI) File Storage. As a cloud engineer, explore file system-specific details through this table, including size, mount targets, and associated metadata. Utilize it to uncover information about file systems, such as their capacity, the associated VCN, and the security of the file system.

## Examples

### Basic info
1. "Explore the basic information about your file systems, such as their display name, ID, lifecycle state, availability domain, metered bytes, and creation time."
2. "Identify instances where file systems are not in an active state, which might indicate potential issues or unused resources."
3. "Discover the file systems that have been cloned, providing a quick overview of your duplicated resources."
4. "Review the file systems that use Oracle managed encryption, to ensure data security and compliance with your organization's encryption policies.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  availability_domain,
  metered_bytes,
  time_created
from
  oci_file_storage_file_system;
```


## List file systems that are not active

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_file_system
where
  lifecycle_state <> 'ACTIVE';
```


## List cloned file systems

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_file_system
where
  is_clone_parent;
```


## List file systems with Oracle managed encryption (default encryption uses Oracle managed encryption keys)

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  time_created
from
  oci_file_storage_file_system
where
  length(kms_key_id) = 0;
```


### List file systems with customer managed encryption keys
Explore file systems that utilize customer-managed encryption keys to enhance data security and privacy. This allows for a better understanding of your system's security measures and helps identify potential areas for improvement.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  kms_key_id,
  time_created
from
  oci_file_storage_file_system
where
  length(kms_key_id) <> 0;
```