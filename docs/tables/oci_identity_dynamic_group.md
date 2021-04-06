# Table: oci_identity_dynamic_group

Dynamic groups allow you to group Oracle Cloud Infrastructure compute instances as "principal" actors (similar to user groups).

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
  oci_identity_dynamic_group;
```


### List dynamic groups which are not in active state

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_identity_dynamic_group
where
  lifecycle_state <> 'ACTIVE';
```


### Get matching rule details of dynamic groups

```sql
select
  name,
  id,
  matching_rule
from
  oci_identity_dynamic_group;
```
