---
title: "Steampipe Table: oci_identity_tag_namespace - Query OCI Identity Tag Namespaces using SQL"
description: "Allows users to query OCI Identity Tag Namespaces."
---

# Table: oci_identity_tag_namespace - Query OCI Identity Tag Namespaces using SQL

Oracle Cloud Infrastructure (OCI) Identity Tag Namespaces is a resource within OCI that assists in categorizing resources in various ways. It allows users to define their own specific set of tags to be used across all resources, helping in better organization and management. This provides a structured way of sorting and filtering resources, thereby improving the efficiency of resource utilization and management.

## Table Usage Guide

The `oci_identity_tag_namespace` table provides insights into the Tag Namespaces within OCI Identity. As a cloud engineer, explore namespace-specific details through this table, including the name, description, and associated metadata. Utilize it to uncover information about namespaces, such as those that are inactive, the compartment they belong to, and the time they were created.

## Examples

### Basic info
Explore which elements within your Oracle Cloud Infrastructure are retired or active. This query provides insights into the lifecycle states of various components, helping you manage your resources effectively.

```sql+postgres
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace;
```

```sql+sqlite
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace;
```

### List active tag namespaces
Determine the areas in which active tag namespaces are present. This is useful for managing and organizing resources within a cloud infrastructure, especially when dealing with a large number of resources.

```sql+postgres
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  lifecycle_state = 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  lifecycle_state = 'ACTIVE';
```

### List retired tag namespaces
Explore which tag namespaces are retired in your Oracle Cloud Infrastructure identity service. This can help maintain an organized and efficient tagging system by identifying and managing outdated tags.

```sql+postgres
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  is_retired;
```

```sql+sqlite
select
  name,
  id,
  is_retired,
  lifecycle_state
from
  oci_identity_tag_namespace
where
  is_retired;
```