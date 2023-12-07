---
title: "Steampipe Table: oci_identity_compartment - Query OCI Identity Compartments using SQL"
description: "Allows users to query OCI Identity Compartments."
---

# Table: oci_identity_compartment - Query OCI Identity Compartments using SQL

Oracle Cloud Infrastructure Identity and Access Management (IAM) service lets you control who has access to your cloud resources. You can control what type of access a group of users have and to which specific resources. This is achieved through the use of compartments, which are a fundamental component of Oracle Cloud Infrastructure for organizing and isolating your cloud resources.

## Table Usage Guide

The `oci_identity_compartment` table provides insights into compartments within Oracle Cloud Infrastructure Identity and Access Management (IAM). As a cloud architect or administrator, you can explore compartment-specific details through this table, including compartment names, descriptions, and states. Utilize it to manage and understand your OCI resource organization, such as identifying compartments with specific resources, understanding your compartment hierarchy, and ensuring appropriate resource isolation.

## Examples

### Basic info
Discover the segments that are in different lifecycle states within your OCI Identity Compartments. This could be useful in assessing the elements within your compartments and understanding their statuses for better resource management.

```sql+postgres
select
  name,
  lifecycle_state,
  id,
  compartment_id as parent_compartment
from
  oci.oci_identity_compartment
```

```sql+sqlite
select
  name,
  lifecycle_state,
  id,
  compartment_id as parent_compartment
from
  oci_identity_compartment
```

### Details joined with parent compartment
Analyze the settings to understand the relationship between different compartments within an OCI identity compartment hierarchy. This is useful for understanding the structure and dependencies within your OCI environment.

```sql+postgres
select
  child.name as name,
  coalesce(parent.name, 'root') as parent_compartment,
  child.id as id,
  coalesce(parent.id, child.tenant_id) as parent_compartment_id
from
  oci_identity_compartment child
  left join oci_identity_compartment parent on (child.compartment_id = parent.id)
order by
  parent.name;
```

```sql+sqlite
select
  child.name as name,
  coalesce(parent.name, 'root') as parent_compartment,
  child.id as id,
  coalesce(parent.id, child.tenant_id) as parent_compartment_id
from
  oci_identity_compartment as child
  left join oci_identity_compartment as parent on child.compartment_id = parent.id
order by
  parent.name;
```

### Full path of the compartments
This query is useful in tracking the full path of active compartments within a system. It aids in system organization and management by allowing users to understand the hierarchical structure of compartments, thereby facilitating easier navigation and data retrieval.

```sql+postgres
with recursive compartments as
(
  select
    name,
    id,
    compartment_id,
    tenant_id,
    name as path,
    name as last_name,
    id as last_id
  from
    oci_identity_compartment
  where
    lifecycle_state = 'ACTIVE'
  union all
  select
    oci_identity_compartment.name,
    oci_identity_compartment.id,
    oci_identity_compartment.compartment_id,
    oci_identity_compartment.tenant_id,
    oci_identity_compartment.name || '/' || compartments.path,
    compartments.last_name,
    compartments.last_id
  from
    oci_identity_compartment
    join
      compartments
      on oci_identity_compartment.id = compartments.compartment_id
)
select
  last_id as compartment_id,
  last_name as name,
  path
from
  compartments
where
  compartment_id = tenant_id
order by
  path;
```

```sql+sqlite
with recursive compartments as
(
  select
    name,
    id,
    compartment_id,
    tenant_id,
    name as path,
    name as last_name,
    id as last_id
  from
    oci_identity_compartment
  where
    lifecycle_state = 'ACTIVE'
  union all
  select
    oci_identity_compartment.name,
    oci_identity_compartment.id,
    oci_identity_compartment.compartment_id,
    oci_identity_compartment.tenant_id,
    oci_identity_compartment.name || '/' || compartments.path,
    compartments.last_name,
    compartments.last_id
  from
    oci_identity_compartment
    join
      compartments
      on oci_identity_compartment.id = compartments.compartment_id
)
select
  last_id as compartment_id,
  last_name as name,
  path
from
  compartments
where
  compartment_id = tenant_id
order by
  path;
```