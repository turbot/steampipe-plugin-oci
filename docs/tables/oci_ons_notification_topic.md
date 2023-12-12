---
title: "Steampipe Table: oci_ons_notification_topic - Query OCI Notification Service Topics using SQL"
description: "Allows users to query OCI Notification Service Topics."
---

# Table: oci_ons_notification_topic - Query OCI Notification Service Topics using SQL

The Oracle Cloud Infrastructure (OCI) Notification Service is a cloud native messaging service that allows for the broadcast of messages to distributed components, such as email and PagerDuty. It is used to send notifications to large numbers of recipients, and ensures that messages are delivered to all subscribers. This service is often used for alerting, coordination, and response in distributed systems.

## Table Usage Guide

The `oci_ons_notification_topic` table provides insights into topics within OCI Notification Service. As a Systems Administrator, explore topic-specific details through this table, including the topic's ARN, name, and status. Utilize it to uncover information about topics, such as those with a specific status, the messages sent through each topic, and the verification of topic configurations.

## Examples

### Basic info
Explore the various notification topics within your Oracle Cloud Infrastructure to understand their lifecycle states and associated API endpoints. This could be useful for assessing the overall configuration and status of your notification system.

```sql+postgres
select
  name,
  topic_id,
  api_endpoint,
  short_topic_id,
  lifecycle_state,
  description
from
  oci_ons_notification_topic;
```

```sql+sqlite
select
  name,
  topic_id,
  api_endpoint,
  short_topic_id,
  lifecycle_state,
  description
from
  oci_ons_notification_topic;
```

### List inactive topics
Explore which notification topics are currently inactive within your Oracle Cloud Infrastructure. This can help you identify areas that may need attention or cleanup to optimize your resources.

```sql+postgres
select
  name,
  lifecycle_state
from
  oci_ons_notification_topic
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  lifecycle_state
from
  oci_ons_notification_topic
where
  lifecycle_state <> 'ACTIVE';
```