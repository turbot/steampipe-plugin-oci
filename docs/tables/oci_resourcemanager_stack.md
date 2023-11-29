---
title: "Steampipe Table: oci_resourcemanager_stack - Query OCI Resource Manager Stacks using SQL"
description: "Allows users to query OCI Resource Manager Stacks."
---

# Table: oci_resourcemanager_stack - Query OCI Resource Manager Stacks using SQL

Oracle Cloud Infrastructure (OCI) Resource Manager is a fully managed service that allows you to automate the process of provisioning your Oracle Cloud Infrastructure resources. Using Terraform, Resource Manager helps you install, configure, and manage resources through the "infrastructure-as-code" model. It provides a consistent, reproducible way to create, change, and improve infrastructure.

## Table Usage Guide

The `oci_resourcemanager_stack` table provides insights into the stacks within OCI Resource Manager. As a DevOps engineer, this table allows you to explore stack-specific details, including configurations, terraform versions, and associated metadata. Utilize it to uncover information about stacks, such as their lifecycle state, time created, and the description of the stack.

## Examples

### Basic info
Explore the status and creation times of your resource manager stacks to understand their lifecycle and manage resources effectively. This can help in assessing the elements within your infrastructure for better planning and resource allocation.

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resourcemanager_stack;
```

### List resource manager stacks that are not active
Determine the areas in which resource manager stacks are not currently active. This can help in identifying unused resources, potentially optimizing resource usage and reducing costs.

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resourcemanager_stack
where
  lifecycle_state <> 'ACTIVE';
```

### List resource manager stacks older than 90 days
Identify instances where resource manager stacks have been in existence for more than 90 days. This can be useful for cleaning up old resources and optimizing resource management.

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_resourcemanager_stack
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```