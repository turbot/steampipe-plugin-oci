# Table: oci_resource_search

OCI resource query search lets you query any and all compartments in the specified tenancy to find resources that match the specified criteria.

**You must specify a Query or Text** in a `where` clause (`where query=' or where text='`).

## Examples

### Freetext search example

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  text = 'database';
```

### List running instances

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query instance resources where lifeCycleState = "RUNNING"';
```

### List resources created in the root compartment

```sql
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query all resources where compartmentId = "ocid1.tenancy.oc1..aaaaaaah5soecxzjetci3yjjnjqmfkr4po3"';
```
