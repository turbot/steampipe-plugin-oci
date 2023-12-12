---
title: "Steampipe Table: oci_core_subnet - Query OCI Core Subnets using SQL"
description: "Allows users to query OCI Core Subnets."
---

# Table: oci_core_subnet - Query OCI Core Subnets using SQL

A Subnet in OCI Core is a subdivision within a VCN (Virtual Cloud Network). It consists of a contiguous range of IP addresses that do not overlap with other subnets in the VCN. You can launch instances, databases, and other resources in a subnet.

## Table Usage Guide

The `oci_core_subnet` table provides insights into subnets within Oracle Cloud Infrastructure's Core Services. As a network administrator, explore subnet-specific details through this table, including associated security lists, route tables, and the availability domain in which the subnet is located. Utilize it to uncover information about subnets, such as those with public IP addresses, the CIDR blocks they cover, and the route tables they use.

## Examples

### Basic info
Explore which subnets are currently active in your Oracle Cloud Infrastructure, along with their creation times and any associated tags. This can help manage resources and track their usage over time.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  tags
from
  oci_core_subnet;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  tags
from
  oci_core_subnet;
```


### Get the OCIDs of the security list for each subnet
Determine the unique identifiers of security lists associated with each network subnet. This can be useful to assess the security configuration and understand how it is applied across different network segments.

```sql+postgres
select
  display_name,
  id,
  jsonb_array_elements_text(security_list_ids) as security_list_id
from
  oci_core_subnet;
```

```sql+sqlite
select
  display_name,
  id,
  json_each.value as security_list_id
from
  oci_core_subnet,
  json_each(security_list_ids);
```


### Count of subnets by VCN ID
Analyze the settings to understand the distribution of subnets across different Virtual Cloud Networks (VCNs). This is useful for managing resources and load balancing across multiple VCNs.

```sql+postgres
select
  vcn_id,
  count(id) as subnet_count
from
  oci_core_subnet
group by
  vcn_id;
```

```sql+sqlite
select
  vcn_id,
  count(id) as subnet_count
from
  oci_core_subnet
group by
  vcn_id;
```


### Get the number of available IP address in each subnet
Explore which subnets have the most available IP addresses to optimize network resource allocation and ensure efficient use of your network space. This is particularly useful in planning network expansion or monitoring network usage.

```sql+postgres
select
  id,
  cidr_block,
  power(2, 32 - masklen(cidr_block :: cidr)) -1 as raw_size
from
  oci_core_subnet;
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```