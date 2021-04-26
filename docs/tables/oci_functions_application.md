# Table: oci_functions_application

In Oracle Functions, an application is:

- a logical grouping of functions
- a way to allocate and configure resources for all functions in the application
- a common context to store configuration variables that are available to all functions in the application
- a way to ensure function runtime isolation

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


### List applications which are not in active state

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


### Get configuration details of each application

```sql
select
  display_name,
  id,
  config
from
  oci_functions_application;
```
