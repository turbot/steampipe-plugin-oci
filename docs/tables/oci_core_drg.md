---
title: "Steampipe Table: oci_core_drg - Query OCI Core Dynamic Routing Gateways using SQL"
description: "Allows users to query Dynamic Routing Gateways in Oracle Cloud Infrastructure's Core Services."
---

# Table: oci_core_drg - Query OCI Core Dynamic Routing Gateways using SQL

Dynamic Routing Gateways (DRGs) are virtual routers that provide a path for private network traffic between your Virtual Cloud Network (VCN) and networks outside the VCN. A DRG is a critical component for creating a site-to-site VPN connection, or a connection that uses Oracle Cloud Infrastructure FastConnect. DRGs provide a secure and reliable connection to your workloads in the Oracle Cloud.

## Table Usage Guide

The `oci_core_drg` table gives insights into Dynamic Routing Gateways within Oracle Cloud Infrastructure's Core Services. As a network administrator, you can delve into details about each DRG, including its state, lifecycle details, and associated compartment. Use this table to manage and monitor your DRGs, ensuring secure and efficient connections between your VCN and external networks.

## Examples

### Basic info
Gain insights into your Oracle Cloud Infrastructure by examining the lifecycle state and creation time of each resource. This can be useful for tracking resource usage and understanding the overall health and status of your infrastructure.

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created
from
  oci_core_drg;
```


### List unavailable dynamic routing gateways
Determine the areas in which dynamic routing gateways are not currently available. This is beneficial to quickly identify and address any network connectivity issues.

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_drg
where
  lifecycle_state <> 'AVAILABLE';
```


### Count of dynamic routing gateways per region
Explore the distribution of dynamic routing gateways across different regions. This can help in managing network traffic and ensuring efficient data routing.

```sql
select
  region,
  count(*) drg_count
from
  oci_core_drg
group by
  region;
```