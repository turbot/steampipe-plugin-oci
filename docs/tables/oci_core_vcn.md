---
title: "Steampipe Table: oci_core_vcn - Query OCI Core Virtual Cloud Networks using SQL"
description: "Allows users to query Virtual Cloud Networks (VCNs) in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_core_vcn - Query OCI Core Virtual Cloud Networks using SQL

A Virtual Cloud Network (VCN) is a customizable, private network in OCI. It closely resembles a traditional network, with firewall rules and specific types of communication gateways that you can choose to use. A VCN resides in a single OCI region and covers a single, contiguous CIDR block of your choice.

## Table Usage Guide

The `oci_core_vcn` table provides insights into Virtual Cloud Networks within Oracle Cloud Infrastructure (OCI). As a network administrator, you can explore network-specific details through this table, including the state of the VCN, its CIDR block, and associated DNS label. Utilize it to uncover information about your VCNs, such as their default security list, route table, and whether they are using DNS resolution.

## Examples

### Basic info
Explore which Virtual Cloud Networks (VCNs) are active and the associated metadata for each. This can be useful in gaining insights into your network's configuration and identifying any potential issues or areas for optimization.

```sql
select
  display_name,
  id,
  lifecycle_state,
  cidr_block,
  freeform_tags
from
  oci_core_vcn;
```

### List unavailable virtual cloud networks
Explore which virtual cloud networks are currently unavailable. This can be useful in identifying potential issues with network connectivity or resource allocation.

```sql
select
  display_name,
  id,
  lifecycle_state
from
  oci_core_vcn
where
  lifecycle_state <> 'AVAILABLE';
```