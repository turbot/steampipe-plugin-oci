# Table: oci_cloud_guard_responder_recipe

A responder is action that Cloud Guard can take when a detector has identified a problem. The available actions are resource-specific. Each responder uses a responder recipe that defines the action or set of actions to take in response to a problem that a detector has identified.

## Examples

### Basic info

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_responder_recipe;
```
