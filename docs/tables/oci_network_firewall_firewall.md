---
title: "Steampipe Table: oci_network_firewall_firewall - Query OCI Network Firewall Firewalls using SQL"
description: "Allows users to query OCI Network Firewall Firewalls."
---

# Table: oci_network_firewall_firewall - Query OCI Network Firewall Firewalls using SQL

The Oracle Cloud Infrastructure (OCI) Network Firewall is a feature of OCI's Virtual Cloud Network (VCN) that provides a security boundary to protect your cloud network resources. It acts as a barrier between your VCN and the public internet, monitoring and controlling incoming and outgoing network traffic based on predetermined security rules. This Firewall service helps in enhancing the security of your cloud resources by reducing the exposure to threats.

## Table Usage Guide

The `oci_network_firewall_firewall` table provides insights into the firewalls within OCI's Network Firewall service. As a Security Analyst, you can explore firewall-specific details through this table, including the associated VCN, the default actions for the firewall's rule sets, and other metadata. Utilize this table to uncover information about your firewalls, such as their current state, the time they were created, and their internet gateway settings.

## Examples

### Basic info
Explore the basic details of your Oracle Cloud Infrastructure network firewalls to gain insights into their availability domains, IP addresses, associated security groups, and current lifecycle state. This allows for efficient management and monitoring of your network security.

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall;
```

### List network firewalls created in the last 30 days
Explore which network firewalls have been created in the past 30 days. This insight can help in assessing recent changes in your network security landscape, enabling you to better manage and monitor your infrastructure's security.

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall
where
  time_created >= now() - interval '30' day;
```

### List network firewalls having IPv6 address
Identify network firewalls that have been assigned an IPv6 address. This can be useful for managing network security and ensuring all devices are properly configured for IPv6 connectivity.

```sql
select
  id,
  display_name,
  availability_domain,
  ipv4_address,
  ipv6_address,
  network_firewall_policy_id,
  network_security_group_ids,
  subnet_id,
  lifecycle_state as state
from
  oci_network_firewall_firewall
where
  ipv6_address is not null;
```

### Describe the network firewall policy associated to the network firewall
Explore the association between network firewalls and their corresponding policies. This can be useful for understanding the lifecycle details of the policy and determining the firewall's adherence to it.

```sql
select
  f.display_name as firewall_name,
  f.id as firewall_id,
  p.display_name as policy_display_name,
  p.id as policy_id,
  p.lifecycle_details as policy_lifecycle 
from
  oci_network_firewall_firewall as f 
  left join
    oci_network_firewall_policy as p 
    on f.network_firewall_policy_id = p.id;
```