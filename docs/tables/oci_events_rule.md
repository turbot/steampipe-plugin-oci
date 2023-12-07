---
title: "Steampipe Table: oci_events_rule - Query OCI Events Rules using SQL"
description: "Allows users to query OCI Events Rules."
---

# Table: oci_events_rule - Query OCI Events Rules using SQL

Oracle Cloud Infrastructure's (OCI) Events service is a cloud-native event routing platform that enables you to automate your application architecture by responding to state changes. Using simple rules, you can match events and route them to one or more targets to take action. A rule is a set of actions that Events service takes when an event that matches a condition occurs.

## Table Usage Guide

The `oci_events_rule` table provides insights into the rules within OCI Events service. As a DevOps engineer, you can explore rule-specific details through this table, including conditions, actions, and associated metadata. Utilize it to uncover information about rules, such as those with specific conditions, the actions taken when those conditions are met, and the verification of those actions.

## Examples

### Basic info
Explore the status and conditions of events rules in Oracle Cloud Infrastructure (OCI) to understand their functionality and lifecycle. This can help in monitoring rule performance, identifying enabled rules, and assessing the actions taken when rules are triggered.

```sql+postgres
select
  id as rule_id,
  display_name,
  lifecycle_state as state,
  condition,
  is_enabled,
  time_created,
  actions
from
  oci_events_rule;
```

```sql+sqlite
select
  id as rule_id,
  display_name,
  lifecycle_state as state,
  condition,
  is_enabled,
  time_created,
  actions
from
  oci_events_rule;
```

### List enabled rules
Explore which rules are currently active in your system. This can help you understand which conditions are being monitored, assisting in system maintenance and troubleshooting.

```sql+postgres
select
  id as rule_id,
  display_name,
  lifecycle_state,
  condition,
  is_enabled
from
  oci_events_rule
where
  is_enabled;
```

```sql+sqlite
select
  id as rule_id,
  display_name,
  lifecycle_state,
  condition,
  is_enabled
from
  oci_events_rule
where
  is_enabled;
```

### List failed rules
Discover the segments that have rules in a 'FAILED' state in your Oracle Cloud Infrastructure's Events service. This can help pinpoint specific areas needing attention, promoting efficient troubleshooting and system optimization.

```sql+postgres
select
  id as rule_id,
  display_name,
  lifecycle_state,
  condition,
  is_enabled
from
  oci_events_rule
where
  lifecycle_state = 'FAILED';
```

```sql+sqlite
select
  id as rule_id,
  display_name,
  lifecycle_state,
  condition,
  is_enabled
from
  oci_events_rule
where
  lifecycle_state = 'FAILED';
```

### Get action details for rules that have the Oracle Notification Service action type
Determine the specifics of rules that utilize the Oracle Notification Service. This query is useful for identifying and managing these particular rules, offering insights into their status, associated topics, and overall configuration.

```sql+postgres
select
  id as rule_id,
  display_name,
  is_enabled,
  a ->> 'actionType' as action_type,
  a ->> 'id' as action_id,
  a ->> 'isEnabled' as action_is_enabled,
  a ->> 'lifecycleState' as action_state,
  a ->> 'topicId' as topic_id
from
  oci_events_rule,
  jsonb_array_elements(actions) as a
where
  a ->> 'actionType'  = 'ONS'
```

```sql+sqlite
select
  rule.id as rule_id,
  display_name,
  is_enabled,
  json_extract(a.value, '$.actionType') as action_type,
  json_extract(a.value, '$.id') as action_id,
  json_extract(a.value, '$.isEnabled') as action_is_enabled,
  json_extract(a.value, '$.lifecycleState') as action_state,
  json_extract(a.value, '$.topicId') as topic_id
from
  oci_events_rule rule,
  json_each(actions) as a
where
  json_extract(a.value, '$.actionType')  = 'ONS'
```

### Get action details for rules that have the Oracle Streaming Service action type
Explore the details of specific rules in Oracle Streaming Service, focusing on their current state and whether they are enabled. This can be useful for managing and monitoring the status of various rules within your streaming service.

```sql+postgres
select
  id as rule_id,
  display_name,
  is_enabled,
  a ->>  'actionType' as action_type,
  a ->> 'id' as action_id,
  a ->> 'isEnabled' as action_is_enabled,
  a ->> 'lifecycleState' as action_state,
  a ->> 'streamId' as stream_id
from
  oci_events_rule,
  jsonb_array_elements(actions) as a
where
  a ->> 'actionType'  = 'OSS'
```

```sql+sqlite
select
  rule_id,
  display_name,
  is_enabled,
  json_extract(a.value, '$.actionType') as action_type,
  json_extract(a.value, '$.id') as action_id,
  json_extract(a.value, '$.isEnabled') as action_is_enabled,
  json_extract(a.value, '$.lifecycleState') as action_state,
  json_extract(a.value, '$.streamId') as stream_id
from
  oci_events_rule,
  json_each(actions) as a
where
  json_extract(a.value, '$.actionType')  = 'OSS'
```

### Get event type details for each rule
Explore which rules are associated with specific event types in your Oracle Cloud Infrastructure events setup. This can help in managing and understanding the functioning of different rules in your system.

```sql+postgres
select
  id as rule_id,
  display_name,
  condition ->> 'eventType' as event_type
from
  oci_events_rule;
```

```sql+sqlite
select
  id as rule_id,
  display_name,
  json_extract(condition, '$.eventType') as event_type
from
  oci_events_rule;
```