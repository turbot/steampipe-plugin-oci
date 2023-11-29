---
title: "Steampipe Table: oci_core_boot_volume_replica - Query OCI Core Boot Volume Replicas using SQL"
description: "Allows users to query data related to Boot Volume Replicas in the Oracle Cloud Infrastructure Core service."
---

# Table: oci_core_boot_volume_replica - Query OCI Core Boot Volume Replicas using SQL

Boot Volume Replicas in the Oracle Cloud Infrastructure Core service are a resource used to create point-in-time, crash-consistent copies of boot volumes. These copies are independent of the source boot volume, allowing you to restore a boot volume to its state at the time the replica was taken. Boot Volume Replicas can be used for disaster recovery, migrating data across regions, and improving backup compliance.

## Table Usage Guide

The `oci_core_boot_volume_replica` table provides insights into Boot Volume Replicas within Oracle Cloud Infrastructure Core service. As a Cloud Engineer, you can use this table to explore details about each replica, including its availability domain, display name, lifecycle state, and associated metadata. This information can be useful for disaster recovery planning, compliance audits, and understanding the state of your boot volume replicas.

## Examples

### Basic info
Explore the basic information of your boot volume replicas to gain insights into their creation time and current state. This can help in assessing their lifecycle management and operational status.

```sql
select
  id,
  display_name,
  boot_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_boot_volume_replica;
```

### List volume replicas which are not available
Discover the segments that consist of volume replicas which are not readily accessible. This can be useful in identifying potential issues or bottlenecks in your storage system.

```sql
select
  id,
  display_name,
  boot_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_boot_volume_replica
where
  lifecycle_state <> 'AVAILABLE';
```