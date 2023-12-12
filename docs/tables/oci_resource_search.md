---
title: "Steampipe Table: oci_resource_search - Query OCI Resource Search using SQL"
description: "Allows users to query OCI Resource Search results."
---

# Table: oci_resource_search - Query OCI Resource Search using SQL

Oracle Cloud Infrastructure (OCI) Resource Search is a service that enables you to find, explore, and understand the entirety of your OCI resources. It provides a unified view of the resources across all compartments in a tenancy, allowing you to manage and govern your resources more effectively. With OCI Resource Search, you can search for resources based on their attributes, relationships, and configurations.

## Table Usage Guide

The `oci_resource_search` table provides insights into the resources within Oracle Cloud Infrastructure. As a cloud administrator, explore resource-specific details through this table, including resource types, compartments, and associated metadata. Utilize it to uncover information about resources, such as their configurations, relationships, and the compartment in which they reside.

**Important Notes**
- You must specify either the `text` or `query` in the `where` clause to query this table.

## Examples

### List resources/services/documentations with freetext as test
Explore which resources, services, or documentations were created at a specific time and are currently active, based on a free text search for the term 'test'. This can help in identifying relevant resources quickly and efficiently.

```sql+postgres
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  text = 'test';
```

```sql+sqlite
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  text = 'test';
```

### List running instances
Identify instances where specific resources are currently active. This allows for a better understanding of resource usage and management, particularly for optimizing operations and troubleshooting.

```sql+postgres
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query instance resources where lifeCycleState = "RUNNING"';
```

```sql+sqlite
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query instance resources where lifeCycleState = "RUNNING"';
```

### List resources created in the root compartment
Determine the areas in which resources were created in the root compartment. This can be beneficial for assessing the elements within your organization's infrastructure, to maintain an overview of resource allocation and usage.

```sql+postgres
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query all resources where compartmentId = "ocid1.tenancy.oc1..aaaaaaah5soecxzjetci3yjjnjqmfkr4po3"';
```

```sql+sqlite
select
  identifier,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resource_search
where
  query = 'query all resources where compartmentId = "ocid1.tenancy.oc1..aaaaaaah5soecxzjetci3yjjnjqmfkr4po3"';
```