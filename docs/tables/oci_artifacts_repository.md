# Table: oci_artifacts_repository

OCI Artifact Repository is a service provided by Oracle Cloud Infrastructure (OCI) that allows you to store and manage various types of artifacts, such as software packages, container images, and other binary files.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  is_immutable,
  description,
  lifecycle_state as state
from
  oci_artifacts_repository;
```

### List immutable repositories

```sql
select
  display_name,
  id,
  time_created,
  is_immutable,
  description,
  lifecycle_state
from
  oci_artifacts_repository
where
  is_immutable;
```

### List repositories created in last 30 days

```sql
select
  display_name,
  id,
  time_created,
  is_immutable,
  description,
  lifecycle_state
from
  oci_artifacts_repository
where
  time_created >= now() - interval '30' day;
```

### List available repositories

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state
from
  oci_artifacts_repository
where
  lifecycle_state = 'AVAILABLE';
```