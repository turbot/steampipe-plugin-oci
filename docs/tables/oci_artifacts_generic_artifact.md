# Table: oci_artifacts_generic_artifact

OCI container artifact information.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  repository_id,
  artifact_path,
  version,
  sha256,
  size_in_bytes,
  lifecycle_state as state
from
  oci_artifacts_generic_artifact;
```