---
title: "Steampipe Table: oci_core_cluster_network - Query OCI Core Cluster Networks using SQL"
description: "Allows users to query information related to cluster networks in Oracle Cloud Infrastructure's Core service."
---

# Table: oci_core_cluster_network - Query OCI Core Cluster Networks using SQL

A Cluster Network is a type of resource provided by the Oracle Cloud Infrastructure's (OCI) Core service. It is a virtual network in the cloud with a scalable, low-latency network infrastructure. Cluster Networks provide a simple, flexible, and secure environment for your compute instances.

## Table Usage Guide

The `oci_core_cluster_network` table provides insights into cluster networks within Oracle Cloud Infrastructure's Core service. As a Network Administrator, you can explore details specific to each cluster network through this table, including its associated instances, security rules, and metadata. Utilize it to uncover information about the network's state, its capacity for instances, and the time it was created.

## Examples

### Basic info
Explore which OCI core cluster networks have been recently created or updated, by identifying their display names and IDs. This can help in understanding their current lifecycle states, which is beneficial for effective network management.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network;
```

### Get instance pool details of cluster network
Analyze the settings to understand the state and size of instance pools within a cluster network. This can help in assessing the overall capacity and availability of the network.

```sql+postgres
select
  c.display_name,
  p -> 'availabilityDomains' as availability_domains,
  p ->> 'instanceConfigurationId' as instance_configuration_id,
  p ->> 'lifecycleState' as instance_pool_state,
  p ->> 'size' as instance_pool_size
from
  oci_core_cluster_network as c,
  jsonb_array_elements(instance_pools) as p;
```

```sql+sqlite
select
  c.display_name,
  json_extract(p.value, '$.availabilityDomains') as availability_domains,
  json_extract(p.value, '$.instanceConfigurationId') as instance_configuration_id,
  json_extract(p.value, '$.lifecycleState') as instance_pool_state,
  json_extract(p.value, '$.size') as instance_pool_size
from
  oci_core_cluster_network as c,
  json_each(c.instance_pools) as p;
```

### List stopped cluster networks
Discover the segments that are associated with halted cluster networks. This is particularly useful for managing resources and ensuring optimal system performance.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  lifecycle_state = 'STOPPED';
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  lifecycle_state = 'STOPPED';
```

### List cluster networks created in the last 30 days
Discover the cluster networks that have been established in the past month. This can be useful for tracking recent network changes and understanding the current lifecycle state of these networks.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_cluster_network
where
  time_created >= datetime('now', '-30 day');
```