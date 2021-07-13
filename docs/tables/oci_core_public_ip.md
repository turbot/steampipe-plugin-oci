# Table: oci_core_public_ip

A public IP address is an IPv4 address that is reachable from the internet. If a resource in your tenancy needs to be directly reachable from the internet, it must have a public IP address. Depending on the type of resource, there might be other requirements.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  ip_address,
  scope,
  time_created
from
  oci_core_public_ip;
```


### List of unused public IPs

```sql
select
  display_name,
  lifecycle_state as state,
  scope
from
  oci_core_public_ip
where
  lifecycle_state = 'AVAILABLE';
```
