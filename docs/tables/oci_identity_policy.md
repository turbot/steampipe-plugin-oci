---
title: "Steampipe Table: oci_identity_policy - Query OCI Identity Policies using SQL"
description: "Allows users to query OCI Identity Policies."
---

# Table: oci_identity_policy - Query OCI Identity Policies using SQL

An OCI Identity Policy is a document that specifies who can access which Oracle Cloud Infrastructure resources that your company has, and the type of access. It is an essential part of Oracle Cloud Infrastructure's Identity and Access Management (IAM) service. It allows you to control who has access to your cloud resources.

## Table Usage Guide

The `oci_identity_policy` table provides insights into Identity Policies within Oracle Cloud Infrastructure's Identity and Access Management (IAM). As a Security Analyst, you can explore policy-specific details through this table, including policy statements, and associated metadata. Utilize it to uncover information about policies, such as those with specific permissions and the verification of policy statements.

## Examples

### Basic info
Explore the lifecycle states and descriptions of various policies within your Oracle Cloud Infrastructure (OCI) to gain insights into their current status and purpose. This can help in managing and understanding your OCI resources effectively.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  description
from
  oci_identity_policy;
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  description
from
  oci_identity_policy;
```

### Get a specific policy
This query allows you to pinpoint the specific details of a chosen policy within your Oracle Cloud Infrastructure (OCI) environment. It's particularly useful when you need to review the configuration or lifecycle state of a policy, without having to sift through all the policies in your OCI setup.

```sql+postgres
select
  name,
  id,
  lifecycle_state,
  description
from
  oci_identity_policy
where
  id = 'ocid1.policy.oc1..aaaaaaaa6nsa2ia2bkr7bx7olpkjuj42yk3rpalwuwvm7fjc7kz7o5wz5pmq';
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state,
  description
from
  oci_identity_policy
where
  id = 'ocid1.policy.oc1..aaaaaaaa6nsa2ia2bkr7bx7olpkjuj42yk3rpalwuwvm7fjc7kz7o5wz5pmq';
```

### List inactive policies
Discover the segments that consist of inactive policies to better manage your resources and maintain an organized, efficient system. This helps in identifying policies that are no longer in use, allowing for effective cleanup and resource optimization.

```sql+postgres
select
  name,
  lifecycle_state
from
  oci_identity_policy
where lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  lifecycle_state
from
  oci_identity_policy
where lifecycle_state <> 'ACTIVE';
```

### List of policy statements
Explore the various policy statements within your Oracle Cloud Infrastructure to better manage and understand your security configuration. This can be particularly useful in identifying potential security loopholes or areas for policy optimization.

```sql+postgres
select
  name,
  jsonb_array_elements_text(statements) as statement
from
  oci_identity_policy
```

```sql+sqlite
select
  name,
  json_each.value as statement
from
  oci_identity_policy,
  json_each(statements)
```