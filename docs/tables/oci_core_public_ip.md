---
title: "Steampipe Table: oci_core_public_ip - Query OCI Core Public IPs using SQL"
description: "Allows users to query details about Public IP resources in the Oracle Cloud Infrastructure Core Services."
---

# Table: oci_core_public_ip - Query OCI Core Public IPs using SQL

A Public IP is a resource in Oracle Cloud Infrastructure Core Services that provides a means to access an instance, a load balancer, or other resource from outside the VCN. The resource can be ephemeral (assigned and unassigned by Oracle) or reserved (you can keep it for as long as you like). Public IPs have scope, which affects their behavior and how they're managed.

## Table Usage Guide

The `oci_core_public_ip` table provides insights into Public IP resources within Oracle Cloud Infrastructure Core Services. As a network administrator or security engineer, explore Public IP-specific details through this table, including their assigned entities, lifecycles, and scopes. Utilize it to uncover information about Public IPs, such as those with reserved lifetimes, the resources they are assigned to, and their current availability status.

## Examples

### Basic info
Analyze the settings to understand the lifecycle state and creation time of your public IP addresses in Oracle Cloud Infrastructure. This can help you manage and track your resources more effectively.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  ip_address,
  scope,
  time_created
from
  oci_core_public_ip;
```

### List unused public IPs
Discover the segments that consist of unused public IPs to optimize resource allocation and reduce unnecessary costs. This is particularly useful in managing network resources effectively by identifying IPs that are available for use.

```sql
select
  display_name,
  lifecycle_state as state,
  scope
from
  oci_core_public_ip
where
  lifecycle_state = 'AVAILABLE';
```