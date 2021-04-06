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

### Get details of specific subscription

```sql
select
  id,
  topic_id,
  lifecycle_state,
  etag,
  freeform_tags,
  delivery_policy
from
  oci_ons_subscription
where
  id = 'ocid1.onssubscription.oc1.ap-mumbai-1.aaaaaaaap57juvblbjgzcddis37gxummh3voqsou54n7eoymkz38uyhbgfdc';
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