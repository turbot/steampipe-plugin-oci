---
title: "Steampipe Table: oci_identity_authentication_policy - Query OCI Identity Authentication Policies using SQL"
description: "Allows users to query OCI Identity Authentication Policies."
---

# Table: oci_identity_authentication_policy - Query OCI Identity Authentication Policies using SQL

Oracle Cloud Infrastructure (OCI) Identity Authentication Policies are a set of rules that define the actions allowed on resources within a compartment. These policies are used to manage access to OCI resources, ensuring that only authorized users can perform operations on these resources. The policies are written in a human-readable, declarative language, and they specify who can access which resources, and how.

## Table Usage Guide

The `oci_identity_authentication_policy` table provides insights into the authentication policies within OCI Identity. As a security administrator, you can explore policy-specific details through this table, including the policy statements, versions, and associated metadata. Utilize it to uncover information about policies, such as those with specific permissions, the resources they affect, and the conditions under which they apply.

## Examples

### Basic info
Uncover the details of your authentication policy to ensure it meets your security standards. This query helps in assessing the password requirements and restrictions, such as length and character type, as well as identifying any network sources that are allowed.

```sql
select
  minimum_password_length,
  is_lowercase_characters_required,
  is_numeric_characters_required,
  is_special_characters_required,
  is_uppercase_characters_required,
  is_username_containment_allowed,
  network_source_ids
from
  oci_identity_authentication_policy
```