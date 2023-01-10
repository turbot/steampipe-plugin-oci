# Table: oci_artifacts_repository

OCI Artifacts repository.

## Examples

### Basic info

```sql
select
    id,
    display_name,
    description,
    is_immutable,
    repository_type,
    lifecycle_state as state
from
    oci_artifacts_repository;
```