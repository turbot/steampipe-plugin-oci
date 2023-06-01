# Table: oci_artifact_container_image

OCI artifact container image is a standardized format for packaging and distributing software applications, libraries, and dependencies. It is based on the Open Container Initiative (OCI) specifications and provides a consistent way to store and share container images across different container runtimes and platforms. OCI artifact container images typically include the application code, runtime environment, system libraries, and other dependencies needed to run the software. They are used to create and deploy containerized applications in cloud environments and container orchestration platforms.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  created_by,
  digest,
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
  oci_artifact_container_image;
```

### Get the largest layers of image in size

```sql
select
  display_name,
  id,
  digest,
  time_created,
  layers_size_in_bytes
from
  oci_artifact_container_image
order by
  layers_size_in_bytes desc limit 1;
```

### Get version details of each images

```sql
select
  i.display_name,
  i.id as image_id,
  v ->> 'createdBy' as image_version_created_by,
  v ->> 'timeCreated' as image_version_created_time,
  v ->> 'version' as version
from
  oci_artifact_container_image as i,
  jsonb_array_elements(versions) as v;
```

### Get layer details of each image

```sql
select
  display_name,
  id,
  l ->> 'digest' as layer_digest,
  l ->> 'sizeInBytes' as layer_size_in_bytes,
  l ->> 'timeCreated' as layer_create_time
from
  oci_artifact_container_image,
  jsonb_array_elements(layers) as l;
```

### Get repository details of each image

```sql
select
  i.display_name,
  i.id,
  i.repository_id,
  r.display_name as repository_display_name,
  r.is_immutable,
  r.is_public,
  r.lifecycle_state
from
  oci_artifact_container_image as i,
  oci_artifact_container_repository as r
where
  r.id = i.repository_id;
```

### List available images

```sql
select
  display_name,
  id,
  digest,
  version,
  lifecycle_state
from
  oci_artifact_container_image
where
  lifecycle_state = 'AVAILABLE';
```

### List images created in last 30 days

```sql
select
  display_name,
  id,
  digest,
  time_created,
  manifest_size_in_bytes
from
  oci_artifact_container_image
where
  time_created >= now() - interval '30' day;
```

### Retrive the total number of pull count for images

```sql
select
  display_name,
  id,
  digest,
  pull_count
from
  oci_artifact_container_image;
```