---
title: "Steampipe Table: oci_core_boot_volume - Query OCI Core Boot Volumes using SQL"
description: "Allows users to query OCI Core Boot Volumes."
---

# Table: oci_core_boot_volume - Query OCI Core Boot Volumes using SQL

A Boot Volume in Oracle Cloud Infrastructure (OCI) Core service is a type of block storage volume that contains the image used to boot a Compute instance. These volumes are reliable and durable with built-in repair capabilities. Boot Volumes are designed for high performance and can be easily backed up, cloned or expanded, as needed.

## Table Usage Guide

The `oci_core_boot_volume` table provides insights into Boot Volumes within OCI Core service. As a cloud administrator, explore Boot Volume-specific details through this table, including its configuration, performance, and associated metadata. Utilize it to uncover information about Boot Volumes, such as their state, size, and the attached instances, to ensure optimal configuration and performance.

## Examples

### Basic info
Explore which boot volumes are in different lifecycle states and when they were created. This can help in managing and tracking the resources effectively.

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume;
```

### List boot volumes with faulty state
Explore which boot volumes are in a faulty state, helping you quickly identify and address potential issues to maintain optimal system performance.

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  lifecycle_state = 'FAULTY';
```

### List boot volumes with size greater than 1024 GB
Explore which boot volumes exceed a size of 1024 GB to manage storage allocation and optimize resource utilization.

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_boot_volume
where
  size_in_gbs > 1024;
```

### List boot volumes with Oracle managed encryption
Discover the segments that have boot volumes managed by Oracle without encryption. This can help identify potential security risks and areas that may need enhanced data protection.
**Note:** Volumes are encrypted by default with Oracle managed encryption key


```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is null;
```

### List boot volumes with customer managed encryption
Gain insights into the boot volumes that have been encrypted using customer-managed keys. This is useful for ensuring compliance with security policies requiring user-managed encryption.

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is not null;
```

### List boot volumes created in the root compartment
Explore which boot volumes have been created in the root compartment of your infrastructure. This can be useful to manage storage resources and ensure efficient use of space.

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_boot_volume
where
  compartment_id = tenant_id;
```