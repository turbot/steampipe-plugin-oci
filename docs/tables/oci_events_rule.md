# Table: oci_events_rule

The Oracle Cloud Infrastructure Events service invokes the action specified in the rule by delivering the event message to action resources, which can include topics, streams, or functions.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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

### Get action details for rules that have the Oracle Streaming Service action type

```sql
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

### Get event type details for each rule

```sql
select
  id as rule_id,
  display_name,
  condition ->> 'eventType' as event_type
from
  oci_events_rule;
```
