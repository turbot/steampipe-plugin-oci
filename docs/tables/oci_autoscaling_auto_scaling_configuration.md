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