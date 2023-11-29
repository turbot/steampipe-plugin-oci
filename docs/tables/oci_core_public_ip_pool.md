---
title: "Steampipe Table: oci_core_public_ip_pool - Query OCI Core Public IP Pools using SQL"
description: "Allows users to query information about Public IP Pools in Oracle Cloud Infrastructure's Core service."
---

# Table: oci_core_public_ip_pool - Query OCI Core Public IP Pools using SQL

Public IP Pools in Oracle Cloud Infrastructure's Core service are resources that represent a pool of public IP addresses. These pools are used to allocate public IPs to resources within the infrastructure. They allow for the management and allocation of public IP addresses in a controlled manner.

## Table Usage Guide

The `oci_core_public_ip_pool` table provides insights into Public IP Pools within Oracle Cloud Infrastructure's Core service. If you are a network administrator or a cloud engineer, you can explore pool-specific details through this table, including the pool's capacity, the number of available IP addresses, and associated metadata. Use it to manage and monitor your public IP address allocation, ensuring optimal use of resources and preventing IP address exhaustion.

## Examples

### Basic info
Explore the lifecycle status, creation time, and region of your Oracle Cloud Infrastructure's public IP pools. This can be useful to understand the distribution and availability of your public IP resources.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  time_created,
  region
from
  oci_core_public_ip_pool;
```

### List public IP pool which are not active
Discover the segments that consist of public IP pools that are not currently active. This can be useful to identify and manage unused resources within your network infrastructure.

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_core_public_ip_pool
where
  lifecycle_state <> 'ACTIVE';
```