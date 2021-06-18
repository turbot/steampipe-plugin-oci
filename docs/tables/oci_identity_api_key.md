# Table: oci_identity_api_key

This resource provides the Api Key resource in Oracle Cloud Infrastructure Identity service.

## Examples

### Basic info

```sql
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key;
```

### List inactive API keys

```sql
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key;
where
  lifecycle_state = 'INACTIVE';
```

### Count API keys by user

```sql
select
  user_id,
  count (*) as api_key_count
from
  oci_identity_api_key
group by
  user_id;
```
