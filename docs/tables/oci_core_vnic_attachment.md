---
title: "Steampipe Table: oci_core_vnic_attachment - Query OCI Core VNIC Attachments using SQL"
description: "Allows users to query VNIC Attachments in OCI Core."
---

# Table: oci_core_vnic_attachment - Query OCI Core VNIC Attachments using SQL

A VNIC attachment represents the attachment of a VNIC to an instance. It is a part of Oracle Cloud Infrastructure's (OCI) Core Services. This resource is used to manage and control the network interface cards (NICs) attached to the instances in your virtual cloud network (VCN).

## Table Usage Guide

The `oci_core_vnic_attachment` table provides insights into VNIC attachments within OCI Core Services. As a Network Administrator, you can explore detailed information about each VNIC attachment, including its lifecycle state, availability domain, and associated instance. Utilize this table to manage and monitor the network interface cards attached to your instances, ensuring optimal network performance and security.

## Examples

### Basic info
Explore which virtual network interface cards (VNICs) are attached to your instances in Oracle Cloud Infrastructure, allowing you to assess their lifecycle states and connectivity details. This information can assist in managing network configurations and troubleshooting network-related issues.

```sql
select
  id,
  display_name,
  instance_id,
  availability_domain,
  lifecycle_state,
  private_ip,
  public_ip,
  time_created
from
  oci_core_vnic_attachment;
```

### List detached VNICs
Explore which VNIC attachments are currently in a 'detached' state. This can be useful for identifying unused resources or potential configuration issues within your OCI environment.

```sql
select
  id,
  display_name,
  lifecycle_state
from
  oci_core_vnic_attachment
where
  lifecycle_state = 'DETACHED';
```

### List automatically created and attached VNICs
Explore which Virtual Network Interface Cards (VNICs) were automatically generated and attached. This is useful for understanding how your network resources are being utilized and managed within your cloud environment.

```sql
select
  id,
  display_name,
  is_primary
from
  oci_core_vnic_attachment
where
  is_primary;
```

### List VNICs with disabled source/destination check
Explore which Virtual Network Interface Cards (VNICs) have their source/destination check disabled. This is useful in situations where you need to route traffic through your instances, such as setting up a NAT (Network Address Translation) instance or configuring a firewall.

```sql
select
  id,
  display_name,
  skip_source_dest_check
from
  oci_core_vnic_attachment
where
  skip_source_dest_check;
```