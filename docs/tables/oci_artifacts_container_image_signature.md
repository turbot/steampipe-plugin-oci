# Table: oci_artifacts_container_image_signature

OCI (Oracle Cloud Infrastructure) Artifact Container Image Signature is a digital signature that provides a way to verify the authenticity and integrity of container images stored in OCI Artifact Registry. It is used to ensure that the image has not been tampered with and that it can be trusted.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  created_by,
  image_id,
  kms_key_id,
  kms_key_version_id,
  message,
  signature,
  signing_algorithm
from
  oci_artifacts_container_image_signature;
```

### List signatures created in last 30 days

```sql
select
  display_name,
  id,
  time_created,
  image_id,
  message,
  signature
from
  oci_artifacts_container_image_signature
where
  time_created >= now() - interval '30' day;
```

### Get image details of each signature

```sql
select
  s.display_name,
  s.id,
  s.signature,
  s.signing_algorithm,
  s.image_id,
  i.digest,
  i.lifecycle_state,
  i.manifest_size_in_bytes,
  i.pull_count
from
  oci_artifacts_container_image_signature as s,
  oci_artifacts_container_image as i
where
  i.id = s.image_id;
```

### Get KMS key details used by each image signature

```sql
select
  s.display_name,
  s.id,
  s.kms_key_version_id,
  v.key_id,
  v.vault_id,
  v.public_key,
  v.origin
from
  oci_artifacts_container_image_signature as s,
  oci_kms_key_version as v
where
  v.id = s.kms_key_version_id;
```

### List signatures with RSA signining algorithm

```sql
select
  display_name,
  id,
  message,
  signature,
  signing_algorithm
from
  oci_artifacts_container_image_signature
where
  signing_algorithm = 'RSA';
```