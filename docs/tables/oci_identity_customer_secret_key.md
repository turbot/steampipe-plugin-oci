# Table: oci_identity_customer_secret_key

A CustomerSecretKey is an Oracle-provided key for using the Object Storage Service's Amazon S3 compatible API. The key consists of a secret key/access key pair. A user can have up to two secret keys at a time.

## Examples

### List of customer secret keys with their corresponding user and date of creation

```sql
select
  id,
  display_name,
  user_id,
  time_created
from
  oci_identity_customer_secret_key;
```


### List customer secret keys which are inactive

```sql
select
  id,
  display_name,
  user_id,
  time_created
from
  oci_identity_customer_secret_key
where
  lifecycle_state = 'INACTIVE';
```


### customer secret keys count per user

```sql
select
  user_id,
  count (id) as customer_secret_key_count
from
  oci_identity_customer_secret_key
group by
  user_id;
```