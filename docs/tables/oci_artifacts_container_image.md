# Table: oci_artifacts_container_image

OCI artifacts container image.

## Examples

### Basic info

```sql
select
    created_by,
    digest,
    display_name,
    id,
    layers,
    layers_size_in_bytes,
    manifest_size_in_bytes,
    pull_count,
    repository_id,
    repository_name,
    versions,
    time_last_pulled,
    version,
    lifecycle_state as state
from
    oci_artifacts_container_image;
```