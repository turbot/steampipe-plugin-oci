---
title: "Steampipe Table: oci_core_network_load_balancer - Query OCI Core Network Load Balancers using SQL"
description: "Allows users to query OCI Core Network Load Balancers."
---

# Table: oci_core_network_load_balancer - Query OCI Core Network Load Balancers using SQL

A Network Load Balancer in Oracle Cloud Infrastructure (OCI) is a regional, non-internet-facing, load balancer that distributes traffic within a virtual cloud network (VCN). It uses a load balancing algorithm and health check policy to distribute traffic among backend servers. The Network Load Balancer is designed to handle volatile traffic patterns and millions of flows, with the ability to scale in real time, without pre-warming.

## Table Usage Guide

The `oci_core_network_load_balancer` table provides insights into Network Load Balancers within Oracle Cloud Infrastructure (OCI) Core services. As a Network Administrator, you can explore load balancer-specific details through this table, including configurations, backend sets, and associated metadata. Utilize it to uncover information about load balancers, such as their health check policies, backend sets, and the distribution of traffic among backend servers.

## Examples

### Basic info
Explore which network load balancers are currently active in your OCI core network. This can help you assess their health status and identify any that may have been recently created or modified.

```sql+postgres
select
  display_name,
  id,
  subnet_id,
  lifecycle_state as state,
  health_status,
  time_created
from
  oci_core_network_load_balancer;
```

```sql+sqlite
select
  display_name,
  id,
  subnet_id,
  lifecycle_state as state,
  health_status,
  time_created
from
  oci_core_network_load_balancer;
```

### List NLBs assigns with public IP address
Discover the segments that are assigned with public IP addresses within your network load balancer, allowing you to identify those that are not private. This can be beneficial for understanding your network's exposure and managing security risks.

```sql+postgres
select
  display_name,
  id,
  is_private
from
  oci_core_network_load_balancer
where
  not is_private;
```

```sql+sqlite
select
  display_name,
  id,
  is_private
from
  oci_core_network_load_balancer
where
  not is_private;
```

### List critical NLBs
Analyze the health status of your network load balancers to identify those in a critical state. This information could be vital in troubleshooting network issues, ensuring data flow efficiency, and maintaining overall system performance.

```sql+postgres
select
  display_name,
  id,
  network_load_balancer_health -> 'status' as health_status
from
  oci_core_network_load_balancer
where
  network_load_balancer_health ->> 'status' = 'CRITICAL';
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(network_load_balancer_health, '$.status') as health_status
from
  oci_core_network_load_balancer
where
  json_extract(network_load_balancer_health, '$.status') = 'CRITICAL';
```