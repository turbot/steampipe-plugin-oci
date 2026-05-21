---
title: "Steampipe Table: oci_waf_network_address_list - Query OCI WAF Network Address Lists using SQL"
description: "Allows users to query OCI WAF Network Address Lists."
---

# Table: oci_waf_network_address_list - Query OCI WAF Network Address Lists using SQL

Oracle Cloud Infrastructure (OCI) Web Application Firewall (WAF) Network Address Lists define sets of IP addresses or VCN address ranges that can be referenced in WAF policies. They support two types: `ADDRESSES` (CIDR-based IP prefixes) and `VCN_ADDRESSES` (private address ranges on an OCI VCN). Network address lists allow you to manage and reuse address sets across multiple WAF policies and rules.

## Table Usage Guide

The `oci_waf_network_address_list` table provides insights into network address lists within OCI WAF. As a security engineer, you can explore address list details through this table, including type, addresses, lifecycle state, and associated metadata. Use it to audit IP allowlists and blocklists, review VCN-scoped address definitions, and track the lifecycle of your WAF network address configurations.

## Examples

### Basic info
Explore the network address lists configured in OCI WAF to understand their type and current lifecycle state.

```sql+postgres
select
  display_name,
  id,
  type,
  lifecycle_state
from
  oci_waf_network_address_list;
```

```sql+sqlite
select
  display_name,
  id,
  type,
  lifecycle_state
from
  oci_waf_network_address_list;
```

### List network address lists not in the active state
Identify network address lists that are not active to detect provisioning issues or resources pending deletion.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_network_address_list
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state
from
  oci_waf_network_address_list
where
  lifecycle_state <> 'ACTIVE';
```

### List CIDR-based address lists with their addresses
Review all `ADDRESSES`-type network address lists and the CIDR prefixes they contain.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(addresses) as addresses
from
  oci_waf_network_address_list
where
  type = 'ADDRESSES';
```

```sql+sqlite
select
  display_name,
  id,
  addresses
from
  oci_waf_network_address_list
where
  type = 'ADDRESSES';
```

### List VCN-based address lists
Identify network address lists that reference private address ranges on an OCI VCN.

```sql+postgres
select
  display_name,
  id,
  jsonb_pretty(vcn_addresses) as vcn_addresses
from
  oci_waf_network_address_list
where
  type = 'VCN_ADDRESSES';
```

```sql+sqlite
select
  display_name,
  id,
  vcn_addresses
from
  oci_waf_network_address_list
where
  type = 'VCN_ADDRESSES';
```

### List network address lists by compartment
Group network address lists by compartment to understand their distribution across your tenancy.

```sql+postgres
select
  compartment_id,
  count(*) as total,
  count(*) filter (where lifecycle_state = 'ACTIVE') as active
from
  oci_waf_network_address_list
group by
  compartment_id;
```

```sql+sqlite
select
  compartment_id,
  count(*) as total,
  sum(case when lifecycle_state = 'ACTIVE' then 1 else 0 end) as active
from
  oci_waf_network_address_list
group by
  compartment_id;
```
