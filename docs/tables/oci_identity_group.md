---
title: "Steampipe Table: oci_identity_group - Query OCI Identity Service Groups using SQL"
description: "Allows users to query OCI Identity Service Groups."
---

# Table: oci_identity_group - Query OCI Identity Service Groups using SQL

Oracle Cloud Infrastructure (OCI) Identity and Access Management (IAM) service lets you control who has access to your cloud resources. You can control what type of access a group of users have and to which specific resources. This is fundamental to the security and compliance of your resources in OCI.

## Table Usage Guide

The `oci_identity_group` table provides insights into the groups within OCI Identity and Access Management (IAM). As a security analyst, you can explore group-specific details through this table, including the users that belong to a group, the policies attached to a group, and other associated metadata. Use it to discover information about groups, such as those with broad permissions, the relationships between users and groups, and to verify policy attachments.

## Examples

### Basic info
Explore which identity groups have been created in your OCI environment, along with their lifecycle states, to understand their current status and when they were established. This could be useful for auditing purposes or for maintaining an overview of your security settings.

```sql
select
  name,
  id,
  description,
  lifecycle_state,
  time_created
from
  oci_identity_group;
```


### List of Identity Groups which are not in Active state
Discover the segments that consist of identity groups not currently in an active state. This is beneficial in identifying and managing inactive groups within your Oracle Cloud Infrastructure.

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_identity_group
where
  lifecycle_state <> 'ACTIVE';
```


### List of Identity Groups without application tag key
Determine the areas in which identity groups lack an application tag key. This is useful for identifying potential gaps in your tagging strategy, helping to ensure all groups are properly categorized and managed.

```sql
select
  name,
  id
from
  oci_identity_group
where
  not tags :: JSONB ? 'application';
```