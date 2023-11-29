---
title: "Steampipe Table: oci_artifacts_repository - Query OCI Artifacts Repositories using SQL"
description: "Allows users to query OCI Artifacts Repositories"
---

# Table: oci_artifacts_repository - Query OCI Artifacts Repositories using SQL

Oracle Cloud Infrastructure (OCI) Artifacts is a fully managed, scalable, and secure service for sharing and tracking software packages. The service can store and distribute packages, manage dependencies, and control versioning. Artifacts supports a variety of package types, including Docker images and Java libraries, among others.

## Table Usage Guide

The `oci_artifacts_repository` table provides insights into repositories within OCI Artifacts. DevOps engineers and developers can gain detailed visibility into their repositories through this table, including their configurations, associated packages, and metadata. Use it to manage and track software packages, control versioning, and ensure the integrity and security of your software supply chain.


## Examples

### Basic info
Explore the basic information of your OCI artifact repositories to understand their immutability, lifecycle state, and other key details. This can be useful in managing your repositories and ensuring proper configuration and state.

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
Explore which repositories have been marked as immutable in your Oracle Cloud Infrastructure. This can be useful for understanding your data's security and compliance status, as immutable repositories cannot be changed or deleted, ensuring the integrity of stored data.

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
Discover the latest repositories that have been created in the last 30 days to understand their lifecycle state and immutability status. This can be useful for auditing purposes, ensuring that all recent additions comply with your organization's policies and standards.

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
Explore which artifact repositories are currently available. This is useful for assessing the resources you have at your disposal for storing and managing your software artifacts.

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