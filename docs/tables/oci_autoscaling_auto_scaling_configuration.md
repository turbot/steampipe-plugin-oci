# Table: oci_autoscaling_auto_scaling_configuration

An autoscaling configuration lets you dynamically scale the resources in a Compute instance pool.

## Examples

### Basic info

```sql
select
    id,
    resource,
    policies,
    display_name,
    cool_down_in_seconds,
    is_enabled,
    max_resource_count,
    min_resource_count
from
    oci_autoscaling_auto_scaling_configuration;
```

### List enabled autoscaling configurations

```sql
select
  display_name,
  id,
  is_enabled
from
  oci_autoscaling_auto_scaling_configuration
where
  is_enabled;
```


### Get policy details for each autoscaling configuration

```sql
select
  display_name as autoscaling_configuration_display_name,
  id,
  p ->> 'displayName' as policy_display_name,
  p ->> 'id' as policy_id,
  p ->> 'isEnabled' as policy_is_enabled,
  p ->> 'policyType' as policy_type,
  p ->> 'rules' as rules,
  p ->> 'capacity' as capacity
from
  oci_autoscaling_auto_scaling_configuration,
  jsonb_array_elements(policies) as p
```


### Get resource details for each autoscaling configuration

```sql
select
  display_name as autoscaling_configuration_display_name,
  id as autoscaling_configuration_id,
  resource ->> 'id' as resource_id,
  resource ->> 'type' as resource_type
from
  oci_autoscaling_auto_scaling_configuration;
```
