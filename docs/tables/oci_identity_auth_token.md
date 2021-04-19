# Table: oci_identity_auth_token

An AuthToken is an Oracle-generated token string that you can use to authenticate with third-party APIs that do not support Oracle Cloud Infrastructure’s signature-based authentication.

## Examples

### List auth tokens with their corresponding user and date of creation

```sql
select
  id,
  user_id,
  user_name,
  time_created
from
  oci_identity_auth_token;
```


### List auth tokens which are inactive

```sql
select
  id,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_auth_token
where
  lifecycle_state = 'INACTIVE';
```


### Auth tokens count per user

```sql
select
  user_id,
  user_name,
  count (id) as auth_token_count
from
  oci_identity_auth_token
group by
  user_name,
  user_id;
```


### List auth tokens older than 90 days

```sql
select
  id,
  user_id,
  user_name,
  lifecycle_state,
  time_created
from
  oci_identity_auth_token
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```