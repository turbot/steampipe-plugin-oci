# Table: oci_devops_repository

Oracle Cloud Infrastructure DevOps repository, is a service provided by Oracle Cloud Infrastructure (OCI) that allows users to securely store, manage, and version control their source code, build artifacts, and Docker images. It is designed to support collaborative development and enable seamless integration with various CI/CD (Continuous Integration/Continuous Deployment) tools and workflows.

## Examples

### Basic info

```sql
select
  id,
  name,
  project_id,
  project_name,
  namespace,
  ssh_url,
  http_url,
  default_branch,
  repository_type,
  time_created,
  lifecycle_state
from
  oci_devops_repository;
```

### List repositories that are not active

```sql
select
  id,
  name,
  project_id,
  project_name,
  repository_type,
  time_created,
  lifecycle_state
from
  oci_devops_repository
where
  lifecycle_state <> 'ACTIVE';
```

### List hosted repositories

```sql
select
  id,
  name,
  project_name,
  repository_type,
  time_created,
  lifecycle_state
from
  oci_devops_repository
where
  repository_type = 'HOSTED';
```

### Count numbers of commits and branches for each repository

```sql
select
  name,
  id,
  branch_count,
  commit_count
from
  oci_devops_repository;
```

### List repositories created in last 30 days

```sql
select
  name,
  id,
  repository_type,
  time_created,
  lifecycle_state
from
  oci_devops_repository
where
  time_created >= now() - interval '30' day;
```