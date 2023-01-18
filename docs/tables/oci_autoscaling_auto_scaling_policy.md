# Table: oci_autoscaling_auto_scaling_policy

Autoscaling policies define the criteria that trigger autoscaling actions and the actions to take. An autoscaling policy is part of an autoscaling configuration.

## Examples

### Basic info

```sql
select
  capacity,
  id,
  display_name,
  is_enabled,
  policy_type 
from
  oci_autoscaling_auto_scaling_policy;
```
