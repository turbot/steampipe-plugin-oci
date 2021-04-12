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
