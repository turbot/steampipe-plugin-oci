# Table: oci_advanced_resource_query_search

OCI advanced resource query search lets you query any and all compartments in the specified tenancy to find resources that match the specified criteria.

## Examples

### Basic info

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_advanced_resource_query_search
where
  query = 'query all resources';
```

### List resources created in the root compartment

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_advanced_resource_query_search
where
  query = 'query all resources where compartmentId = "ocid1.tenancy.oc1..aaaaaaah5soecxzjetci3yjjnjqmfkr4po3"';
```

### List running instances

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_advanced_resource_query_search
where
  query = 'query instance resources where lifeCycleState = "RUNNING"';
```
