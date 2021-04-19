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


### List NSG which are not in availabe state

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


### Count of NSG per vcn

```sql
select
  vcn_id,
  count(id) as no_of_nsg
from
  oci_core_network_security_group
group by
  vcn_id;
```