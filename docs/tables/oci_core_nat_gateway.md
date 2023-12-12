---
title: "Steampipe Table: oci_core_nat_gateway - Query OCI Core NAT Gateways using SQL"
description: "Allows users to query NAT Gateways in the Oracle Cloud Infrastructure Core Service."
---

# Table: oci_core_nat_gateway - Query OCI Core NAT Gateways using SQL

A NAT Gateway in the Oracle Cloud Infrastructure Core Service is a networking component that enables instances in a private subnet to connect to the internet or other Oracle services, but prevents the internet from initiating connections with those instances. It provides a path for private network resources to access the internet, with the NAT Gateway acting as a public-facing, internet-resolvable IP address. The NAT Gateway ensures that return traffic is directed to the correct resource.

## Table Usage Guide

The `oci_core_nat_gateway` table provides insights into the NAT Gateways within the Oracle Cloud Infrastructure Core Service. As a network engineer or cloud architect, explore NAT Gateway-specific details through this table, including its lifecycle state, block traffic IP address, and associated network details. Utilize it to uncover information about NAT Gateways, such as their associated subnets, the public IP address used, and the verification of their operational status.

## Examples

### Basic info
Explore which NAT gateways are currently active within your Oracle Cloud Infrastructure. This query can help you assess the state and region of each gateway, as well as any associated tags, providing a comprehensive overview of your network's security and routing configuration.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  public_ip_id,
  region,
  tags
from
  oci_core_nat_gateway;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  public_ip_id,
  region,
  tags
from
  oci_core_nat_gateway;
```


### List NAT Gateways that blocks traffic
Determine the areas in which NAT Gateways are blocking traffic to assess potential network bottlenecks or security measures.

```sql+postgres
select
  display_name,
  id,
  block_traffic
from
  oci_core_nat_gateway
where
  block_traffic;
```

```sql+sqlite
select
  display_name,
  id,
  block_traffic
from
  oci_core_nat_gateway
where
  block_traffic = 1;
```


### Count NAT gateways by VCN
Assess the distribution of NAT gateways across your virtual cloud networks to better understand your network infrastructure and optimize resource allocation.

```sql+postgres
select
  vcn_id,
  count(*) as nat_gateway_count
from
  oci_core_nat_gateway
group by
  vcn_id;
```

```sql+sqlite
select
  vcn_id,
  count(*) as nat_gateway_count
from
  oci_core_nat_gateway
group by
  vcn_id;
```