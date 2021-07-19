# Table: oci_core_network_load_balancer

A network security group (NSG) provides a virtual firewall for a set of cloud resources that all have the same security posture.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  subnet_id,
  lifecycle_state as state,
  health_status,
  time_created
from
  oci_core_network_load_balancer;
```


### List NLBs assigns with public IP address

```sql
select
  display_name,
  id,
  is_private
from
  oci_core_network_load_balancer
where
  not is_private;
```


### List unhealthy NLBs

```sql
select
  display_name,
  id,
  health_status
from
  oci_core_network_load_balancer
where
  health_status <> 'OK';
```
