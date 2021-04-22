# Table: oci_identity_customer_secret_key

A customer secret key is an Oracle-provided key for using the Object Storage Service's Amazon S3 compatible API. The key consists of a secret key/access key pair. A user can have up to two secret keys at a time.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  user_id,
  user_name,
  time_created
from
  oci_identity_customer_secret_key;
```


### List inactive customer secret keys

```sql
select
  id,
  display_name,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_customer_secret_key
where
  lifecycle_state = 'INACTIVE';
```


### Count customer secret keys by user

```sql
select
  user_id,
  count (id) as customer_secret_key_count
from
  oci_identity_customer_secret_key
group by
  user_id;
```
