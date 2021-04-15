# Table: oci_ons_subscription

Subscriptions allow you to be notified of new messages or changes via a Message Queue of your choice.

## Examples

### Basic info

```sql
select
  id,
  topic_id,
  lifecycle_state,
  protocol,
  endpoint,
  etag
from
  oci_ons_subscription;
```


### List of subscription where state is pending

```sql
select
  id,
  lifecycle_state,
  protocol,
  endpoint
from
  oci_ons_subscription
where
  lifecycle_state = 'PENDING';
```