# Table: oci_artifacts_container_image_signature

OCI artifacts container image signature.

## Examples

### Basic info

```sql
select
  created_by,
  display_name,
  id,
  image_id,
  kms_key_id,
  kms_key_version_id,
  message,
  signature,
  signing_algorithm
from
  oci_artifacts_container_image_signature;
```