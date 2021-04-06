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
  rr ->> 'cidrBlock' as cidr_block,
  rr ->> 'description' as description,
  rr ->> 'destination' as destination,
  rr ->> 'destinationType' as destination_type,
  rr ->> 'networkEntityId' as network_entity_id
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rr;
```


### List route tables whose routes are directed to the internet

```sql
select
  display_name,
  id,
  rr ->> 'destination' as destination
from
  oci_core_route_table,
  jsonb_array_elements(route_rules) as rr
where
  rr ->> 'destination' = '0.0.0.0/0'
```