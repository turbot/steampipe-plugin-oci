# Table: oci_artifact_container_repository

OCI Artifacts Container Repositories are a service provided by Oracle Cloud Infrastructure (OCI) for storing and managing container images. These repositories allow you to securely store, share, and deploy container images within OCI. They provide a central location to store and manage your container images, enabling you to easily distribute and deploy them across your OCI environment. OCI Artifacts Container Repositories support popular container image formats and can be integrated with other OCI services such as Oracle Container Engine for Kubernetes (OKE) for seamless container image deployment.

## Examples

### Basic info

```sql
select
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
  oci_artifact_container_repository;
```

### List repositories that are not public

```sql
select
  display_name,
  id,
  image_count,
  is_immutable,
  is_public,
  layer_count
from
  oci_artifact_container_repository
where
  not is_public;
```

### List repositories that are immutable

```sql
select
  display_name,
  id,
  is_public,
  is_immutable,
  layer_count,
  layers_size_in_bytes
from
  oci_artifact_container_repository
where
  is_immutable;
```

### List repositories created in last 30 days

```sql
select
  display_name,
  id,
  time_created,
  created_by,
  lifecycle_state
from
  oci_artifact_container_repository
where
  time_created >= now() - interval '30' day;
```

### Get layer details of repositories

```sql
select
  display_name,
  layer_count,
  layers_size_in_bytes
from
  oci_artifact_container_repository;
```

### List top 5 billable repositories

```sql
select
  display_name,
  is_immutable,
  is_public,
  billable_size_in_g_bs
from
  oci_artifact_container_repository
order by
  billable_size_in_g_bs desc limit 5;
```

### List available repositories

```sql
select
  display_name,
  id,
  time_created,
  created_by,
  lifecycle_state
from
  oci_artifact_container_repository
where
  lifecycle_state = 'AVAILABLE';
```