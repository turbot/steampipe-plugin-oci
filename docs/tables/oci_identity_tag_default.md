# Table: oci_identity_tag_default

Tag defaults let you specify a default tag (tagnamespace.tag="value") to apply to all resource types in a specified compartment. The tag default is applied at the time the resource is created. Resources that exist in the compartment before you create the tag default are not tagged. The TagDefault object specifies the tag and compartment details.

## Examples

### Basic info

```sql
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_default;
```

### List active tag namespaces

```sql
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  lifecycle_state = 'ACTIVE';
```

### List retired tag namespaces

```sql
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  is_retired;
```
