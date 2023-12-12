---
title: "Steampipe Table: oci_core_internet_gateway - Query OCI Core Internet Gateways using SQL"
description: "Allows users to query OCI Core Internet Gateways."
---

# Table: oci_core_internet_gateway - Query OCI Core Internet Gateways using SQL

An OCI Core Internet Gateway is a virtual router that provides a path for traffic between your VCN and the internet. As a stateless gateway, it does not maintain information about the traffic that flows through it. It also does not perform Network Address Translation (NAT) or route table lookups.

## Table Usage Guide

The `oci_core_internet_gateway` table provides insights into the Internet Gateways within Oracle Cloud Infrastructure (OCI). As a network engineer, explore gateway-specific details through this table, including its lifecycle state, availability domain, and associated metadata. Utilize it to uncover information about gateways, such as their attached VCNs, the time they were created, and whether they are enabled for internet traffic.

## Examples

### Basic info
Gain insights into the creation and lifecycle status of your internet gateway resources in Oracle Cloud Infrastructure. This is useful for managing and monitoring the state and duration of these resources in your network.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_internet_gateway;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_core_internet_gateway;
```


### List disabled internet gateways
Identify instances where internet gateways are disabled to understand potential gaps in your network infrastructure. This can help in improving the security and reliability of your systems.

```sql+postgres
select
  display_name,
  id,
  time_created,
  is_enabled
from
  oci_core_internet_gateway
where
  not is_enabled;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  is_enabled
from
  oci_core_internet_gateway
where
  not is_enabled;
```