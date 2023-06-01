# Table: oci_artifact_generic_artifact

OCI container artifact information.

## Examples

### Basic info

```sql
select
  id,
  name,
  repository_id,
  artifact_path,
  version,
  sha256,
  size_in_bytes,
  lifecycle_state as state
from
  oci_artifact_generic_artifact;
```

### List available generic artifacts

```sql
select
  name,
  id,
  repository_id,
  artifact_path,
  version,
  size_in_bytes,
  lifecycle_state
from
  oci_artifact_generic_artifact
where
  lifecycle_state = 'AVAILABLE';
```

### Count numbers of artifact versions by artifact path

```sql
select
  artifact_path,
  count(version) as numbers_of_versions
from
  oci_artifact_generic_artifact
group by
  artifact_path;
```

### List generic artifacts created in the last 30 days

```sql
select
  name,
  id,
  sha256,
  lifecycle_state,
  artifact_path,
  time_created
from
  oci_artifact_generic_artifact
where
  time_created >= now() - interval '30' day;
```

### Get the largest artifact

```sql
select
  name,
  id,
  artifact_path,
  version,
  size_in_bytes
from
  oci_artifact_generic_artifact
order by
  size_in_bytes desc limit 1;
```

### Get repository details for an artifact

```sql
select
  a.id,
  a.name as artifact_name,
  r.display_name as repository_display_name,
  r.is_immutable as is_repository_immutable,
  r.lifecycle_state as repository_lifecycle_state
from
  oci_artifact_generic_artifact as a,
  oci_artifact_repository as r
where
  a.repository_id = r.id
and
  a.id = 'ocid1.genericartifact.oc1.ap-mumbai-1.0.amaaaaaa6igdexaaxzyuikdquye6wozpb4rxgkijxe77pfu64zigyqp7o5ua';
```