---
title: "Steampipe Table: oci_artifacts_generic_artifact - Query OCI Artifacts Generic Artifacts using SQL"
description: "Allows users to query OCI Artifacts Generic Artifacts."
---

# Table: oci_artifacts_generic_artifact - Query OCI Artifacts Generic Artifacts using SQL

Oracle Cloud Infrastructure Artifacts service provides a highly scalable and distributed artifact repository. It is designed to store and share container images and other software development artifacts, such as JAR files, Python packages, and so on. This service is integrated with Oracle Cloud Infrastructure Registry, which is a Docker v2 compliant registry and supports Docker CLI and Docker Engine for pushing and pulling images.

## Table Usage Guide

The `oci_artifacts_generic_artifact` table provides insights into generic artifacts within Oracle Cloud Infrastructure Artifacts service. As a software developer or DevOps engineer, explore artifact-specific details through this table, including version, size, and associated metadata. Utilize it to uncover information about artifacts, such as those with specific versions, the state of artifacts, and the verification of artifact metadata.

## Examples

### Basic info
Analyze the settings to understand the state and size of your artifacts in Oracle Cloud Infrastructure. This allows for better management and organization of your resources.

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
  oci_artifacts_generic_artifact;
```

### List available generic artifacts
Explore which generic artifacts are currently available. This can be beneficial in managing resources or tracking the lifecycle of various artifacts within your OCI environment.

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
  oci_artifacts_generic_artifact
where
  lifecycle_state = 'AVAILABLE';
```

### Count numbers of artifact versions by artifact path
Analyze the settings to understand the distribution of versions across different artifact paths. This can help you identify areas where versions are proliferating, potentially indicating a need for version management or cleanup.

```sql
select
  artifact_path,
  count(version) as numbers_of_versions
from
  oci_artifacts_generic_artifact
group by
  artifact_path;
```

### List generic artifacts created in the last 30 days
Discover the recently created generic artifacts within the past month. This is useful for tracking recent activities and changes in your environment.

```sql
select
  name,
  id,
  sha256,
  lifecycle_state,
  artifact_path,
  time_created
from
  oci_artifacts_generic_artifact
where
  time_created >= now() - interval '30' day;
```

### Get the largest artifact
Discover the largest artifact within your OCI environment, which can help you manage storage and identify potential data bottlenecks. This is particularly useful for optimizing storage allocation and improving overall system performance.

```sql
select
  name,
  id,
  artifact_path,
  version,
  size_in_bytes
from
  oci_artifacts_generic_artifact
order by
  size_in_bytes desc limit 1;
```

### Get repository details for an artifact
Gain insights into the characteristics of a specific artifact by analyzing its associated repository details. This can be particularly useful when you need to understand the repository's immutability and lifecycle state for better artifact management.

```sql
select
  a.id,
  a.name as artifact_name,
  r.display_name as repository_display_name,
  r.is_immutable as is_repository_immutable,
  r.lifecycle_state as repository_lifecycle_state
from
  oci_artifacts_generic_artifact as a,
  oci_artifacts_repository as r
where
  a.repository_id = r.id
and
  a.id = 'ocid1.genericartifact.oc1.ap-mumbai-1.0.amaaaaaa6igdexaaxzyuikdquye6wozpb4rxgkijxe77pfu64zigyqp7o5ua';
```