---
title: "Steampipe Table: oci_core_route_table - Query OCI Core Route Tables using SQL"
description: "Allows users to query OCI Core Route Tables."
---

# Table: oci_core_route_table - Query OCI Core Route Tables using SQL

The OCI Core Route Table is a virtual table in Oracle Cloud Infrastructure (OCI) that enables network traffic routing within a Virtual Cloud Network (VCN). It contains rules that determine where network traffic is directed. Each subnet in a VCN uses a route table to direct outgoing traffic.

## Table Usage Guide

The `oci_core_route_table` table provides insights into the route tables within Oracle Cloud Infrastructure's Core Networking service. As a network engineer, you can explore route table-specific details through this table, including associated rules, route targets, and associated metadata. Utilize it to uncover information about route tables, such as those with specific rules, the targets of each route, and the verification of route targets.

## Examples

### Basic info
Explore which routes are active within your network by identifying their creation time and regional distribution. This can help understand the network's structure and management, and optimize its performance.

```sql+postgres
select
  display_name,
  id,
  vcn_id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_route_table;
```

```sql+sqlite
select
  display_name,
  id,
  vcn_id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_route_table;
```


### Get routing details for each route table
Explore the routing details for each network route, gaining insights into their destination types and associated network entities. This is useful in identifying potential network bottlenecks and understanding the network's architecture.

```sql+postgres
select
  display_name,
  id,
  rt ->> 'cidrBlock' as cidr_block,
  rt ->> 'description' as description,
  rt ->> 'destination' as destination,
  rt ->> 'destinationType' as destination_type,
  rt ->> 'networkEntityId' as network_entity_id
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rt;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(rt.value, '$.cidrBlock') as cidr_block,
  json_extract(rt.value, '$.description') as description,
  json_extract(rt.value, '$.destination') as destination,
  json_extract(rt.value, '$.destinationType') as destination_type,
  json_extract(rt.value, '$.networkEntityId') as network_entity_id
from
  oci_core_route_table,
  json_each(route_rules) as rt;
```


### List route tables with routes directed to the Internet
Explore which route tables have routes directed to the Internet. This is particularly helpful for assessing potential security risks and ensuring correct data routing configurations.

```sql+postgres
select
  display_name,
  id,
  rt ->> 'destination' as destination
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rt
where
  rt ->> 'destination' = '0.0.0.0/0'
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(rt.value, '$.destination') as destination
from
  oci_core_route_table,
  json_each(route_rules) as rt
where
  json_extract(rt.value, '$.destination') = '0.0.0.0/0'
```