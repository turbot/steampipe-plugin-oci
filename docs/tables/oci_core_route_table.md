# Table: oci_core_route_table

A collection of RouteRule objects, which are used to route packets based on destination IP to a particular network entity

## Examples

### Basic info

```sql
select
  display_name,
  id,
  vcn_id,
  time_created,
  lifecycle_state as state,
  region
from
  oci_core_route_table;
```


### Get routing details of route table

```sql
select
  display_name,
  id,
  rt ->> 'cidrBlock' as cidr_block,
  rt ->> 'description' as description,
  rt ->> 'destination' as destination,
  rt ->> 'destinationType' as destination_type,
  rt ->> 'networkEntityId' as network_entity_id
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rt;
```


### List route tables whose routes are directed to the internet

```sql
select
  display_name,
  id,
  rt ->> 'destination' as destination
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rt
where
  rt ->> 'destination' = '0.0.0.0/0'
```