---
title: "Steampipe Table: oci_ons_subscription - Query OCI Notification Service Subscriptions using SQL"
description: "Allows users to query OCI Notification Service Subscriptions."
---

# Table: oci_ons_subscription - Query OCI Notification Service Subscriptions using SQL

The Oracle Cloud Infrastructure (OCI) Notification Service is a cloud messaging service for sending notifications to large numbers of recipients. It enables you to broadcast messages to distributed components by topic, through a publish-subscribe pattern. This service ensures that messages are delivered to all active, healthy subscribers in a topic.

## Table Usage Guide

The `oci_ons_subscription` table provides insights into subscriptions within the OCI Notification Service. As a DevOps engineer, explore subscription-specific details through this table, including subscription protocol, endpoint, and associated metadata. Utilize it to uncover information about subscriptions, such as their delivery policy, effective delivery policy, and the status of the subscription.

## Examples

### Basic info
Explore which topics are subscribed to in the Oracle Cloud Infrastructure (OCI) messaging service, understanding the lifecycle state and protocol of each. This aids in monitoring the health and status of your subscribed topics for efficient message processing.

```sql+postgres
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

```sql+sqlite
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
Explore which subscriptions are in a pending state to assess their status and take necessary actions. This is useful in managing subscriptions and ensuring timely activation or termination.

```sql+postgres
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

```sql+sqlite
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
Explore the distribution of subscriptions across different topics. This can be useful to understand which topics are attracting the most subscribers, aiding in content strategy and resource allocation.

```sql+postgres
select
  topic_id,
  count(id) as subscription_count
from
  oci_ons_subscription
group by
  topic_id;
```

```sql+sqlite
select
  topic_id,
  count(id) as subscription_count
from
  oci_ons_subscription
group by
  topic_id;
```

### Get subscription count by protocol
Analyze the distribution of different protocols to understand the overall usage patterns in your Oracle Cloud Infrastructure (OCI) Notification Service subscriptions. This can help in optimizing resource allocation and planning for future requirements.

```sql+postgres
select
  protocol,
  count(protocol) as protocol_count
from
  oci_ons_subscription
group by
  protocol;
```

```sql+sqlite
select
  protocol,
  count(protocol) as protocol_count
from
  oci_ons_subscription
group by
  protocol;
```