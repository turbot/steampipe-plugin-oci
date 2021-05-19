# Table: oci_core_network_security_group

A network security group (NSG) provides a virtual firewall for a set of cloud resources that all have the same security posture.

## Examples

### Basic info

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

```sql
select
  vcn_id,
  count(id) as no_of_nsg
from
  oci_core_network_security_group
group by
  vcn_id;
```
