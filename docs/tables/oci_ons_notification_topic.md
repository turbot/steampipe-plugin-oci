# Table: oci_ons_notification_topic

A topic is a communication channel for sending messages to its subscriptions. A topic can have zero, one, or multiple subscriptions that are notified whenever a message is published to a topic.

## Examples

### Basic info

```sql
select
  name,
  topic_id,
  lifecycle_state,
  description
from
  oci_ons_notification_topic;
```


### List inactive topic

```sql
select
  name,
  lifecycle_state
from
  oci_ons_notification_topic
where lifecycle_state <> 'ACTIVE';
```
