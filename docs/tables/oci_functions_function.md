# Table: oci_functions_function

In Oracle Functions, a function is small but powerful blocks of code that generally do one simple thing. It is grouped into applications, stored as Docker images in a specified Docker registry and invoked in response to a CLI command or signed HTTP request.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  application_id,
  lifecycle_state,
  image,
  image_digest
from
  oci_functions_function;
```


### List functions where trace configuration is disabled

```sql
select
  display_name,
  id,
  application_id,
  trace_config -> 'isEnabled' as trace_config_is_enabled
from
  oci_functions_function
where
  not (trace_config -> 'isEnabled') :: bool;
```


### List functions where memory is greater than 100 MB

```sql
select
  display_name,
  id,
  application_id,
  memory_in_mbs
from
  oci_functions_function
where
  memory_in_mbs > 100;
```
