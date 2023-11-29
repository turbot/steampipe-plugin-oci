---
title: "Steampipe Table: oci_core_local_peering_gateway - Query OCI Core Local Peering Gateways using SQL"
description: "Allows users to query OCI Core Local Peering Gateways."
---

# Table: oci_core_local_peering_gateway - Query OCI Core Local Peering Gateways using SQL

A Local Peering Gateway (LPG) in Oracle Cloud Infrastructure (OCI) is a regional object that represents a local VCN peering relationship. It allows the peering of VCNs in the same region, enabling the VCNs to function as a single, unified network without routing traffic over the internet or through a router. The LPG is a key component of the VCN peering architecture within OCI, facilitating the exchange of private network traffic between VCNs.

## Table Usage Guide

The `oci_core_local_peering_gateway` table provides insights into Local Peering Gateways within Oracle Cloud Infrastructure's Core Services. As a network administrator, explore gateway-specific details through this table, including associated VCNs, lifecycle states, and associated metadata. Utilize it to uncover information about peering gateways, such as those associated with specific VCNs, the lifecycle state of each gateway, and the verification of peering relationships.

## Examples

### Basic info
Explore the status of local peering gateways within your Oracle Cloud Infrastructure network to manage and optimize connections between virtual cloud networks. This is vital for ensuring smooth data transfer and maintaining network performance.

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway;
```

### List available LPGs
Discover the segments that have local peering gateways in an available state, enabling you to manage and optimize network connectivity within your virtual cloud network.

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway
where
  lifecycle_state = 'AVAILABLE';
```

### List LPGs which are not connected to any peer
Explore which Local Peering Gateways (LPGs) are not connected to any peer. This can help in identifying isolated network segments, which could be a potential security risk or a misconfiguration.

```sql
select
  name,
  id,
  vcn_id,
  lifecycle_state
from
  oci_core_local_peering_gateway
where
  peering_status = 'NEW';
```