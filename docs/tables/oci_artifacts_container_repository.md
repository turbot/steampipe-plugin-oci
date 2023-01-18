# Table: oci_artifacts_container_repository

OCI artifacts container repositories.

## Examples

### Basic info

```sql
select
  created_by,
  display_name,
  id,
  image_count,
  is_immutable,
  is_public,
  layer_count,
  layers_size_in_bytes,
  billable_size_in_g_bs,
  time_last_pushed,
  lifecycle_state as state 
from
  oci_artifacts_container_repository;
```
