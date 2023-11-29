---
title: "Steampipe Table: oci_core_network_security_group - Query OCI Core Network Security Groups using SQL"
description: "Allows users to query Network Security Groups in OCI Core service."
---

# Table: oci_core_network_security_group - Query OCI Core Network Security Groups using SQL

A Network Security Group in OCI Core service is a virtual firewall for your virtual network resources. It provides inbound and outbound security rules that specify the type of traffic and the ports on which the traffic is allowed. Network Security Groups are associated with a Virtual Cloud Network (VCN) and can be associated with specific subnets or a specific set of Compute instances.

## Table Usage Guide

The `oci_core_network_security_group` table provides insights into the Network Security Groups within OCI Core service. As a security analyst, you can explore group-specific details through this table, including security rules, associated resources, and other metadata. Leverage this table to identify potential security risks, such as overly permissive security rules and to ensure compliance with your organization's security policies.

## Examples

### Basic info
Explore which network security groups are active within your system and when they were created. This can help maintain security standards and identify any potentially unauthorized or outdated groups.

```sql
select
  display_name,
  id,
  vcn_id,
  lifecycle_state as state,
  time_created
from
  oci_core_network_security_group;
```


### List NSGs that are not available
Discover the segments that are not currently available in your network security groups. This is useful to assess the elements within your network that might be causing issues due to their unavailability.

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_core_network_security_group
where
 lifecycle_state <> 'AVAILABLE';
```


### Count the number of rules per NSG
1. "Explore the distribution of rules across different Network Security Groups (NSG) in your environment."
2. "Identify Network Security Groups (NSGs) in your environment that have unrestricted inbound access from the internet."
3. "Uncover Network Security Groups (NSGs) in your environment that have unrestricted SSH and RDP access from the internet.

```sql
select
  display_name,
  id,
  jsonb_array_length(rules) as rules_count
from
  oci_core_network_security_group;
```


## List NSGs whose inbound access is open to the internet

```sql
select
  display_name,
  id,
  r ->> 'direction' as direction,
  r ->> 'sourceType' as source_type,
  r ->> 'source' as source
from
  oci_core_network_security_group,
  jsonb_array_elements(rules) as r
where
  r ->> 'direction' = 'INGRESS'
  and r ->> 'sourceType' = 'CIDR_BLOCK'
  and r ->> 'source' = '0.0.0.0/0'
```


## List NSG whose SSH and RDP access is not restricted from the internet

```sql
select
  display_name,
  id,
  r ->> 'direction' as direction,
  r ->> 'sourceType' as source_type,
  r ->> 'source' as source,
  r ->> 'protocol' as protocol,
  r -> 'tcpOptions' -> 'destinationPortRange' ->> 'max' as min_port_range,
  r -> 'tcpOptions' -> 'destinationPortRange' ->> 'min' as max_port_range
from
  oci_core_network_security_group,
  jsonb_array_elements(rules) as r
where
  r ->> 'direction' = 'INGRESS'
  and r ->> 'sourceType' = 'CIDR_BLOCK'
  and r ->> 'source' = '0.0.0.0/0'
  and (
    (
      r ->> 'protocol' = 'all'
    )
    or (
      (r -> 'tcpOptions' -> 'destinationPortRange' ->> 'min')::integer <= 22
      and (r -> 'tcpOptions' -> 'destinationPortRange' ->> 'max')::integer >= 22
    )
    or (
      (r -> 'tcpOptions' -> 'destinationPortRange' ->> 'min')::integer <= 3389
      and (r -> 'tcpOptions' -> 'destinationPortRange' ->> 'max')::integer >= 3389
    )
  );
```


### Count the number of NSGs per VCN
Explore which Virtual Cloud Networks (VCNs) have the most Network Security Groups (NSGs) to understand your cloud network's security distribution. This can help in assessing the security coverage and identifying areas that might need additional security measures.

```sql
select
  vcn_id,
  count(id) as no_of_nsg
from
  oci_core_network_security_group
group by
  vcn_id;
```