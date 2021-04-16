# Table: oci_core_security_list

Security lists act as virtual firewalls for your Compute instances and other kinds of resources. A security list consists of a set of ingress and egress security rules that apply to all the VNICs in any subnet that the security list is associated with.

## Examples

### Basic info

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


### Details of egress security rules

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

### Details of ingress security rules

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


### Count of security list by VCN ID

```sql
select
 vcn_id,
 count(id) as security_list_count
from
  oci_core_security_list
group by
  vcn_id;
```


### List of default security list

```sql
select
  display_name,
  id
from
  oci_core_security_list
where
  display_name like '%Default Security%';
```