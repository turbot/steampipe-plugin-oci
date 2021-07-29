# Table: oci_core_public_ip_pool

A public IP pool is a set of public IP addresses represented as one or more IPv4 CIDR blocks. Resources like load balancers and compute instances can be allocated public IP addresses from a public IP pool.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  time_created,
  region
from
  oci_core_public_ip_pool;
```

### List public IP pool which are not active

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_core_public_ip_pool
where
  lifecycle_state <> 'ACTIVE';
```