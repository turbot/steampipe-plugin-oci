---
title: "Steampipe Table: oci_devops_repository - Query OCI DevOps Repositories using SQL"
description: "Allows users to query OCI DevOps Repositories"
---

# Table: oci_devops_repository - Query OCI DevOps Repositories using SQL

OCI DevOps Repositories is a feature of Oracle Cloud Infrastructure's DevOps service that provides a place for storing, sharing, and versioning application source code. It supports Git-based repositories, enabling teams to collaborate and manage their codebase efficiently. This feature is integral to the CI/CD pipeline, facilitating code changes tracking, branching, and merging.

## Table Usage Guide

The `oci_devops_repository` table provides insights into repositories within OCI DevOps service. As a DevOps engineer, you can explore repository-specific details through this table, including repository names, ids, compartment ids, and associated metadata. Utilize it to manage and audit your repositories, understand their configuration, and identify any potential issues or improvements.

## Examples

### Basic info
Explore the basic details of your DevOps repositories, such as identity, associated project, URLs, and status. This can help you manage and understand your projects, especially in large or complex environments.

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
Discover the segments that consist of inactive repositories in your DevOps project. This could be useful to identify and manage unused resources, thereby optimizing your project's performance and cost.

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
Explore which hosted repositories are present in your project by identifying their key details such as ID, name, and creation time. This can be beneficial for understanding the types and states of repositories within your project, aiding in efficient project management.

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
Analyze the activity within each repository by assessing the number of commits and branches. This can provide insights into the level of development and collaboration occurring within each repository, aiding in project management and resource allocation.

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
Discover the segments that were created in the last month. This is useful for tracking recent activity and identifying new additions to your system.

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