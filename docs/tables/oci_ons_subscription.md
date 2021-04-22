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

### List subscriptions in a pending state

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

### Get subscription count by topic ID

```sql
select
  topic_id,
  count(id) as subscription_count
from
  oci_ons_subscription
group by
  topic_id;
```

### Get subscription count by protocol

```sql
select
  protocol,
  count(protocol) as protocol_count
from
  oci_ons_subscription
group by
  protocol;
```
