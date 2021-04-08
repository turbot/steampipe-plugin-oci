# Table: oci_core_service_gateway

A service gateway enables cloud resources without public IP addresses to privately access Oracle services.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  vcn_id lifecycle_state,
  region
from
  oci_core_service_gateway;
```

### List Route tables that Service gateways use

```sql
select
  display_name,
  id,
  route_table_id
from
  oci_core_service_gateway;
```

### Services info that enabled for Service gateways

```sql
select
  display_name,
  id,
  s -> 'serviceId' as service_id,
  s -> 'serviceName' as service_name
from
  oci_core_service_gateway,
  jsonb_array_elements(services) as s;
```


### List Service gateways that block traffics

```sql
select
  display_name,
  id,
  block_traffic
from
  oci_core_service_gateway
where
  block_traffic;
```


### List of Service gateways without application tag key

```sql
select
  display_name,
  id,
  tags
from
  oci_core_service_gateway
where
  not tags :: JSONB ? 'application';
```