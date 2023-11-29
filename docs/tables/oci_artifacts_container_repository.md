---
title: "Steampipe Table: oci_artifacts_container_repository - Query OCI Artifacts Container Repositories using SQL"
description: "Allows users to query OCI Artifacts Container Repositories."
---

# Table: oci_artifacts_container_repository - Query OCI Artifacts Container Repositories using SQL

Oracle Cloud Infrastructure (OCI) Artifacts is a fully managed, scalable, and secure artifact storage and sharing service. It provides a repository for storing and sharing artifacts like Docker images and helps you manage their deployment. Artifacts Container Repository is a part of OCI Artifacts that stores Docker images.

## Table Usage Guide

The `oci_artifacts_container_repository` table provides insights into the OCI Artifacts Container Repositories. DevOps engineers, developers, and system administrators can leverage this table to explore and manage Docker images stored in the repositories. It can be used to uncover information about Docker images, such as their tags, layers, and associated metadata.

## Examples

### Basic info
Explore the basic information of your OCI artifacts container repository to understand its current state and usage. This can be useful in managing resources, assessing security, and tracking changes over time.

```sql
select
  display_name,
  id,
  image_count,
  is_immutable,
  is_public,
  layer_count,
  layers_size_in_bytes,
  billable_size_in_gbs,
  time_last_pushed,
  lifecycle_state as state
from
  oci_artifacts_container_repository;
```

### List repositories that are not public
Discover the segments that contain private repositories in your Oracle Cloud Infrastructure's artifact container. This can be beneficial to assess the elements within your infrastructure that are not publicly accessible, ensuring data privacy and security.

```sql
select
  display_name,
  id,
  image_count,
  is_immutable,
  is_public,
  layer_count
from
  oci_artifacts_container_repository
where
  not is_public;
```

### List repositories that are immutable
Discover the segments that are marked as immutable within a container repository. This is useful for understanding storage details and identifying repositories that cannot be changed, assisting in maintaining data integrity and security.

```sql
select
  display_name,
  id,
  is_public,
  is_immutable,
  layer_count,
  layers_size_in_bytes
from
  oci_artifacts_container_repository
where
  is_immutable;
```

### List repositories created in last 30 days
Identify newly created repositories within the past month. This is useful for tracking recent additions and understanding the growth and activity within your system.

```sql
select
  display_name,
  id,
  time_created,
  created_by,
  lifecycle_state
from
  oci_artifacts_container_repository
where
  time_created >= now() - interval '30' day;
```

### Get layer details of repositories
Explore the layer details of your repositories to understand their size and count. This can be useful for managing storage and optimizing repository performance.

```sql
select
  display_name,
  layer_count,
  layers_size_in_bytes
from
  oci_artifacts_container_repository;
```

### List top 5 billable repositories
Uncover the details of your most resource-intensive repositories. This query helps you identify and prioritize the top five repositories based on their billable size for efficient resource management.

```sql
select
  display_name,
  is_immutable,
  is_public,
  billable_size_in_gbs
from
  oci_artifacts_container_repository
order by
  billable_size_in_gbs desc limit 5;
```

### List available repositories
Discover the segments that are readily accessible in your system by exploring the available repositories. This query is particularly useful in instances where you need to manage the lifecycle of your resources and understand who created them and when.

```sql
select
  display_name,
  id,
  time_created,
  created_by,
  lifecycle_state
from
  oci_artifacts_container_repository
where
  lifecycle_state = 'AVAILABLE';
```