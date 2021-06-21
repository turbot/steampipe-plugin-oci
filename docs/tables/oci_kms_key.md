# Table: oci_kms_key

Oracle Cloud Infrastructure Key Management is a managed service that enables you to manage and control AES symmetric keys used to encrypt your data-at-rest.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  time_created,
  vault_name
from
  oci_kms_key;
```

### List keys which are not enabled

```sql
select
  id,
  name,
  lifecycle_state,
  vault_name
from
  oci_kms_key
where
  lifecycle_state <> 'ENABLED';
```

### List keys older than 365 days

```sql
select
  id,
  name,
  lifecycle_state,
  vault_name
from
  oci_kms_key
where
  time_created <= (current_date - interval '365' day)
order by
  time_created;
```
