---
title: "Steampipe Table: oci_identity_tag_default - Query OCI Identity Tag Defaults using SQL"
description: "Allows users to query OCI Identity Tag Defaults."
---

# Table: oci_identity_tag_default - Query OCI Identity Tag Defaults using SQL

Oracle Cloud Infrastructure (OCI) Identity Tag Defaults is a feature within the OCI Identity service that allows you to define default tags to be applied to all resources at the time of creation within a specific compartment. These default tags help in organizing and tracking resources, improving cost tracking, security, and compliance across the OCI resources. The tag defaults are inherited by all new resources created in the compartment and can be overwritten during or after the resource creation.

## Table Usage Guide

The `oci_identity_tag_default` table provides insights into the default tags within OCI Identity. As a cloud administrator or security analyst, explore tag-specific details through this table, including tag namespace, value, and lifecycle state. Utilize it to uncover information about default tags, such as their applicability to specific resources, their lifecycle states, and the compartments they are associated with.

## Examples

### Basic info
Analyze the settings to understand the necessity and lifecycle status of various elements within your OCI Identity Tag Default. This can help in assessing the importance and current stage of each element, aiding in efficient resource management.

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default;
```

### List active tag defaults
Explore which tag defaults are currently active within your OCI identity configuration. This can help manage and organize your resources effectively.

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default
where
  lifecycle_state = 'ACTIVE';
```

### List required tag defaults
Determine the required tag defaults within your OCI identity to ensure compliance and manage resources effectively. This query is particularly useful in identifying and understanding the lifecycle state of these mandatory tag defaults.

```sql
select
  id,
  is_required,
  lifecycle_state
from
  oci_identity_tag_default
where
  is_required;
```