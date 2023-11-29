---
title: "Steampipe Table: oci_core_volume - Query OCI Core Volumes using SQL"
description: "Allows users to query OCI Core Volumes."
---

# Table: oci_core_volume - Query OCI Core Volumes using SQL

OCI Core Volumes are block storage devices that you can attach to instances in the same availability domain. They are durable, high performance storage volumes that can persist beyond the life of a single Oracle Cloud Infrastructure Compute instance. Core Volumes are designed for mission-critical applications that require granular, consistent, and reliable data access.

## Table Usage Guide

The `oci_core_volume` table provides insights into Core Volumes within Oracle Cloud Infrastructure (OCI). As a database administrator, explore volume-specific details through this table, including size, state, and associated metadata. Utilize it to uncover information about volumes, such as those with high performance, the availability of volumes, and the verification of volume backups.

## Examples

### Basic info
Explore which volumes in your Oracle Cloud Infrastructure are currently active, when they were created, and their display names for easy identification and management.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume;
```


### List volumes with a faulty state
Explore which storage volumes are in a faulty state. This is useful for identifying potential issues in your storage infrastructure that may require troubleshooting or maintenance.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume
where
  lifecycle_state = 'FAULTY';
```


### List volumes with size greater than 1024 GB
Identify instances where storage volumes exceed 1024 GB to effectively manage resources and ensure optimal utilization.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_volume
where
  size_in_gbs > 1024;
```


### List volumes with Oracle managed encryption (volumes are encrypted by default with Oracled managed encryption keys)
Determine the volumes that are encrypted using Oracle's default management keys. This can be useful to identify volumes that may not have been individually configured with custom encryption settings, potentially highlighting areas for enhanced security measures.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume
where
  kms_key_id is null;
```


### List volumes with customer managed encryption
Identify instances where customer-managed encryption is used for volumes. This can help assess security measures and compliance with data protection standards.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume
where
  kms_key_id is not null;
```


### List volumes created in the root compartment
Discover the segments that are utilizing storage space within the root compartment. This can help in efficient management of resources and planning for future storage needs.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_volume
where
  compartment_id = tenant_id;
```