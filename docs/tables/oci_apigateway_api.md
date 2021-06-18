# Table: oci_apigateway_api

The API Gateway service enables you to publish APIs with private endpoints that are accessible from within your network, and which you can expose with public IP addresses if you want them to accept internet traffic.       

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


### List active APIs

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


### List APIs older than 90 days

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
