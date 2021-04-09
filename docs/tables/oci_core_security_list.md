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

### Get a specific security list

```sql
select
  display_name,
  id,
  lifecycle_state,
  egress_security_rules,
  ingress_security_rules
from
  oci_core_security_list
where
  id = 'ocid1.securitylist.oc1.ap-mumbai-1.aaaaaaaak6haggimwbffytlr72267dwzgv3pbps2t467winmu4erfg5wsdcv';
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