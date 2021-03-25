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

### Groups details to which user belongs

```sql
select
  name as user_name,
  groups ->> 'groupId' as group_id
from
  oci_identity_user,
  jsonb_array_elements(user_groups) as groups
```
