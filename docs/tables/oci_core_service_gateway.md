# Table: oci_core_service_gateway

A service gateway enables cloud resources without public IP addresses to privately access Oracle services.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  vcn_id,
  lifecycle_state,
  region
from
  oci_core_service_gateway;
```

### List service gateways that use route tables

```sql
select
  display_name,
  id,
  route_table_id
from
  oci_core_service_gateway
where
  route_table_id is not null;
```

### Get enabled services for each service gateway

```sql
select
  display_name,
  id,
  s ->> 'serviceId' as service_id,
  s ->> 'serviceName' as service_name
from
  oci_core_service_gateway,
  jsonb_array_elements(services) as s;
```

### List service gateways that block traffic

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
