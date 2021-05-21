# Table: oci_cloud_guard_configuration

Settings are the current set of configurations for Cloud Guard. Cloud Guard can only be disabled from the reporting region.

## Examples

### Basic info

```sql
select
  reporting_region,
  status,
  self_manage_resources
from
  oci_cloud_guard_configuration;
```
