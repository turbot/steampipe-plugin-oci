---
title: "Steampipe Table: oci_core_block_volume_replica - Query OCI Core Block Volume Replicas using SQL"
description: "Allows users to query information about block volume replicas in Oracle Cloud Infrastructure's (OCI) core services."
---

# Table: oci_core_block_volume_replica - Query OCI Core Block Volume Replicas using SQL

Block Volume Replicas in Oracle Cloud Infrastructure (OCI) are copies of block volumes that can be used to create new block volumes or replace existing ones. These replicas provide a point-in-time copy of the data and are useful for backup, restore, and cloning operations. They are an integral part of OCI's core services, providing data durability, reliability, and flexibility.

## Table Usage Guide

The `oci_core_block_volume_replica` table provides insights into block volume replicas within OCI's core services. As a database administrator or storage engineer, you can explore replica-specific details through this table, including its availability domain, lifecycle state, and associated metadata. Utilize it to manage your block volume replicas, such as understanding their current state, ensuring data consistency, and planning for disaster recovery scenarios.

## Examples

### Basic info
Explore which block volume replicas are currently active and when they were created. This can help in monitoring and managing the lifecycle of your storage resources.

```sql+postgres
select
  id,
  display_name,
  block_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_block_volume_replica;
```

```sql+sqlite
select
  id,
  display_name,
  block_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_block_volume_replica;
```

### List volume replicas which are not available
Explore the block volume replicas in your OCI Core that are not currently available. This can help in identifying potential issues or disruptions in your data storage and backup systems.

```sql+postgres
select
  id,
  display_name,
  block_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_block_volume_replica
where
  lifecycle_state <> 'AVAILABLE';
```

```sql+sqlite
select
  id,
  display_name,
  block_volume_id,
  time_created,
  lifecycle_state as state
from
  oci_core_block_volume_replica
where
  lifecycle_state <> 'AVAILABLE';
```