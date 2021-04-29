# Table: oci_functions_application

In Oracle Functions, an application is both a unit of function runtime isolation and a logical grouping of related functions. It provides a context to store network configuration and environment variables that are available to all functions in the application.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  lifecycle_state,
  subnet_ids
from
  oci_functions_application;
```


### List applications not in the active state

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_functions_application
where
  lifecycle_state <> 'ACTIVE';
```


### Get configuration details for each application

```sql
select
  display_name,
  id,
  config
from
  oci_functions_application;
```
