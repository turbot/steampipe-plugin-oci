# Table: oci_identity_group

Group is collection of users who all need the same type of access to a particular set of resources or compartment.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  lifecycle_state,
  time_created
from
  oci_identity_group;
```


### List of Identity Groups which are not in Active state

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_identity_group
where
  lifecycle_state <> 'ACTIVE';
```


### List of Identity Groups without application tag key

```sql
select
  name,
  id
from
  oci_identity_group
where
  not tags :: JSONB ? 'application';
```
