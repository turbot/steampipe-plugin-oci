# Table: oci_identity_tag_default

Tag defaults let you specify a default tag (tagnamespace.tag="value") to apply to all resource types in a specified compartment. The tag default is applied at the time the resource is created. Resources that exist in the compartment before you create the tag default are not tagged. The TagDefault object specifies the tag and compartment details.

## Examples

### Basic info

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default;
```

### List active tag defaults

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default
where
  lifecycle_state = 'ACTIVE';
```

### List required tag defaults

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default
where
  is_required;
```
