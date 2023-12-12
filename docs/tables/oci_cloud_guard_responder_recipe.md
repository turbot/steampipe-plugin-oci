---
title: "Steampipe Table: oci_cloud_guard_responder_recipe - Query OCI Cloud Guard Responder Recipes using SQL"
description: "Allows users to query Cloud Guard Responder Recipes in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_cloud_guard_responder_recipe - Query OCI Cloud Guard Responder Recipes using SQL

Cloud Guard Responder Recipe is a component of Oracle Cloud Infrastructure's (OCI) Cloud Guard service, which provides a robust security and compliance platform. Responder Recipes are collections of managed lists and rules that define how Cloud Guard should respond to specific security findings. They play a crucial role in automating the remediation of detected security issues.

## Table Usage Guide

The `oci_cloud_guard_responder_recipe` table provides insights into the Responder Recipes within OCI's Cloud Guard service. As a security analyst or a cloud administrator, explore detailed information about these recipes, including their associated rules, conditions, and actions. Utilize this table to streamline the management of your security responses, automate remediation processes, and enhance the overall security posture of your OCI environment.

## Examples

### Basic info
Explore which cloud guard responder recipes have been created and their current lifecycle states. This is useful for understanding the security posture and state of your Oracle Cloud Infrastructure.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_responder_recipe;
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_responder_recipe;
```

### List responder recipes which are not active
Discover the segments that consist of responder recipes which are not currently active. This can be useful in identifying potential security risks or areas in need of maintenance within your cloud infrastructure.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_responder_recipe
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_responder_recipe
where
  lifecycle_state <> 'ACTIVE';
```