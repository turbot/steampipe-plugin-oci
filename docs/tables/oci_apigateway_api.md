# Table: oci_apigateway_api

API gateway plays an important role as a secure access point that protects an organization's APIs.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_apigateway_api;
```


### List active Api in a compartment

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_apigateway_api
where
  lifecycle_state = 'ACTIVE';
```


### List api older than 90 days

```sql
select
  id,
  lifecycle_state,
  time_created
from
  oci_apigateway_api
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```