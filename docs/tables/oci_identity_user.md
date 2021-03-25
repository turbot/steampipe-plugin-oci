# Table: oci_identity_user

## Examples

### Basic info

```sql
select
  name,
  id,
  user_type,
  time_created,
  is_mfa_activated,
  can_use_console_password
from
  oci_identity_user;
```

### Details of identity groups attached to users

```sql
select
  oci_identity_user.name as user_name,
  oci_identity_group.name as group_name,
  user_group ->> 'groupId' as group_id
from
  oci_identity_user,
  jsonb_array_elements(user_groups) as user_group
  inner join oci_identity_group ON (oci_identity_group.id = user_group ->> 'groupId' )
```
