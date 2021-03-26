# Table: oci_identity_user

## Examples

### Basic info

```sql
select
  name,
  id,
  email,
  user_type,
  time_created,
  lifecycle_state,
  is_mfa_activated,
  can_use_api_keys,
  can_use_console_password,
  can_use_auth_tokens,
  can_use_smtp_credentials,
  can_use_customer_secret_keys
from
  oci_identity_user;
```

### List Oracle Identity Cloud Service(IDCS) users

```sql
select
  name,
  id,
  email,
  time_created,
  lifecycle_state,
  is_mfa_activated
from
  oci_identity_user
where
  user_type = 'IDCS';
```

### List users who can log in to console

```sql
select
  name,
  user_type
from
  oci_identity_user
where
  can_use_console_password;
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
  inner join oci_identity_group ON (oci_identity_group.id = user_group ->> 'groupId' );
```
