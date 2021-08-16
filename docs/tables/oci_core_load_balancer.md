# Table: oci_core_load_balancer

Oracle Cloud Infrastructure (OCI) Flexible Load Balancing enables customers to distribute web requests across a fleet of servers or automatically route traffic across fault domains, availability domains, or regions—yielding high availability and fault tolerance for any application or data source.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  shape_name
from
  oci_core_load_balancer;
```

### List load balancers assigns with public IP address

```sql
select
  display_name,
  id,
  is_private
from
  oci_core_load_balancer
where
  not is_private;
```
