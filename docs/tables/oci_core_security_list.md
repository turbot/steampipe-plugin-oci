---
title: "Steampipe Table: oci_core_security_list - Query OCI Core Services Security Lists using SQL"
description: "Allows users to query Security Lists within OCI Core Services."
---

# Table: oci_core_security_list - Query OCI Core Services Security Lists using SQL

A Security List in OCI Core Services is a virtual firewall for an instance, with ingress and egress rules that specify the types of traffic allowed in and out. Each rule can have a type, destination, protocol, and source. Security Lists are stateful, meaning their rules apply to both incoming and outgoing traffic.

## Table Usage Guide

The `oci_core_security_list` table provides insights into Security Lists within OCI Core Services. As a network administrator, you can examine the details of these lists, including the types of traffic they allow and their rules. This table is a valuable tool for understanding and managing the flow of traffic in and out of your instances.

## Examples

### Basic info
Assess the elements within your network by identifying the lifecycle state and creation time of security lists in Oracle Cloud Infrastructure. This helps pinpoint specific locations where security measures are active, aiding in overall network management and security.

```sql
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  vcn_id
from
  oci_core_security_list;
```


### Get egress security rules for each security list
Uncover the details of egress security rules for each security list to understand their settings and configurations. This can be used to assess the elements within each list, providing insights into the security protocols and options in place.

```sql
select
  display_name,
  p ->> 'destination' as destination,
  p ->> 'destinationType' as destination_type,
  p ->> 'icmpOptions' as icmp_options,
  p ->> 'isStateless' as is_stateless,
  p ->> 'protocol' as protocol,
  p ->> 'tcpOptions' as tcp_options,
  p ->> 'udpOptions' as udp_options
from
  oci_core_security_list,
  jsonb_array_elements(egress_security_rules) as p;
```

### Get ingress security rules for each security list
Determine the areas in which your system's security could be improved by analyzing the ingress security rules for each security list. This allows you to identify potential vulnerabilities and take necessary action to enhance your system's security.

```sql
select
  display_name,
  p ->> 'description' as description,
  p ->> 'icmpOptions' as icmp_options,
  p ->> 'isStateless' as is_stateless,
  p ->> 'protocol' as protocol,
  p ->> 'source' as source,
  p ->> 'sourceType' as source_type,
  p ->> 'tcpOptions' as tcp_options,
  p ->> 'udpOptions' as udp_options
from
  oci_core_security_list,
  jsonb_array_elements(ingress_security_rules) as p;
```


### List security lists that do not restrict SSH and RDP access from the internet
Identify instances where the security lists are not restricting SSH and RDP access from the internet, which could potentially expose your network to security risks.

```sql
 select
  display_name,
  p ->> 'description' as description,
  p ->> 'icmpOptions' as icmp_options,
  p ->> 'isStateless' as is_stateless,
  p ->> 'protocol' as protocol,
  p ->> 'source' as source,
  p ->> 'sourceType' as source_type,
  p -> 'tcpOptions' -> 'destinationPortRange' ->> 'max' as min_port_range,
  p -> 'tcpOptions' -> 'destinationPortRange' ->> 'min' as max_port_range,
  p ->> 'udpOptions' as udp_options
from
  oci_core_security_list,
  jsonb_array_elements(ingress_security_rules) as p
where
  p ->> 'source' = '0.0.0.0/0'
  and (
    (
      p ->> 'protocol' = 'all'
      and (p -> 'tcpOptions' -> 'destinationPortRange' -> 'min') is null
    )
    or (
      (p -> 'tcpOptions' -> 'destinationPortRange' ->> 'min')::integer <= 22
      and (p -> 'tcpOptions' -> 'destinationPortRange' ->> 'max')::integer >= 22
    )
    or (
      (p -> 'tcpOptions' -> 'destinationPortRange' ->> 'min')::integer <= 3389
      and (p -> 'tcpOptions' -> 'destinationPortRange' ->> 'max')::integer >= 3389
    )
  );
```


### List default security lists
Explore the default security lists within your system to understand their unique identifiers and names. This is useful in assessing the existing security configurations and identifying any potential areas of concern.

```sql
select
  display_name,
  id
from
  oci_core_security_list
where
  display_name like '%Default Security%';
```