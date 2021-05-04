# Table: oci_identity_tag_namespace

OCI Tag namespace is a container for your tag keys. It consists of a name and zero or more tag key definitions. Tag namespaces are not case sensitive and must be unique across the tenancy. The namespace is also a natural grouping to which administrators can apply policy. One policy on the tag namespace applies to all the tag definitions contained within that namespace.

## Examples

### Basic info

```sql
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace;
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

### List tag namespaces which are retired

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
