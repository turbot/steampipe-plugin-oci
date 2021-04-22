# Table: oci_cloud_guard_detector_recipe

OCI Cloud Guard detector recipe uses multiple detector rules, each of which a specific definition of a class of resources, with specific actions or configurations, that cause a detector to report a problem. If any one rule is triggered, it causes the detector to report a problem.

## Examples

### Basic info

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe;
```

### List detector recipe which are active

```sql
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe
where
  lifecycle_state = 'ACTIVE';
```

### List detector recipes where password related rules are enabled or disabled

```sql
select
  name,
  e ->> 'detectorRuleId' as Rule_name,
  e -> 'details' ->> 'isEnabled' as status
from
  oci_cloud_guard_detector_recipe,
  jsonb_array_elements(effective_detector_rules) as e
where
  e ->> 'detectorRuleId' = 'PASSWORD_TOO_OLD'
  or e ->> 'detectorRuleId' = 'PASSWORD_POLICY_NOT_COMPLEX';
```
