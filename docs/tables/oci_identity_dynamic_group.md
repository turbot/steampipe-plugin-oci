---
title: "Steampipe Table: oci_identity_dynamic_group - Query OCI Identity Dynamic Groups using SQL"
description: "Allows users to query OCI Identity Dynamic Groups."
---

# Table: oci_identity_dynamic_group - Query OCI Identity Dynamic Groups using SQL

OCI Identity Dynamic Groups are a collection of compute instances within an OCI compartment that match rules defined in a statement. This allows you to manage instances in a dynamic manner without the need to manage each instance individually. It is a key component in OCI's Identity and Access Management service.

## Table Usage Guide

The `oci_identity_dynamic_group` table provides insights into dynamic groups within OCI Identity and Access Management (IAM). As a cloud administrator, explore group-specific details through this table, including group rules, associated compartment ID, and associated metadata. Utilize it to manage and monitor dynamic groups, such as those with specific rules, the relationships between groups, and the verification of group policies.

## Examples

### Basic info
Explore the basic information of dynamic groups in your Oracle Cloud Infrastructure. This can help in understanding their current lifecycle state and when they were created, which is useful for management and auditing purposes.

```sql
select
  name,
  id,
  description,
  lifecycle_state,
  time_created
from
  oci_identity_dynamic_group;
```


### List inactive dynamic groups
Analyze the settings to understand which dynamic groups in your OCI identity are not currently active. This can help manage resources and identify potential security risks.

```sql
select
  name,
  id,
  lifecycle_state
from
  oci_identity_dynamic_group
where
  lifecycle_state <> 'ACTIVE';
```