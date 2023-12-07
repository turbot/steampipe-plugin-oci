---
title: "Steampipe Table: oci_identity_user - Query OCI Identity Users using SQL"
description: "Allows users to query information about users in OCI Identity."
---

# Table: oci_identity_user - Query OCI Identity Users using SQL

Oracle Cloud Infrastructure's Identity and Access Management (IAM) service lets you control who has access to your cloud resources. You can control what type of access a group of users have and to which specific resources. This is done through the use of policies, compartments, and other security features that the IAM service offers.

## Table Usage Guide

The `oci_identity_user` table provides insights into users within OCI Identity and Access Management (IAM). As a security administrator, explore user-specific details through this table, including user ID, name, description, and associated metadata. Utilize it to uncover information about users, such as their state, time of creation, and compartment ID.

## Examples

### Basic info
Discover the segments that highlight user details and their access privileges. This allows for better management and oversight of user permissions and security settings.

```sql+postgres
select
  name,
  id,
  email,
  user_type,
  time_created,
  lifecycle_state,
  is_mfa_activated,
  can_use_api_keys,
  can_use_console_password,
  can_use_auth_tokens,
  can_use_smtp_credentials,
  can_use_customer_secret_keys
from
  oci_identity_user;
```

```sql+sqlite
select
  name,
  id,
  email,
  user_type,
  time_created,
  lifecycle_state,
  is_mfa_activated,
  can_use_api_keys,
  can_use_console_password,
  can_use_auth_tokens,
  can_use_smtp_credentials,
  can_use_customer_secret_keys
from
  oci_identity_user;
```

### List Oracle Identity Cloud Service(IDCS) users
Explore which users in the Oracle Identity Cloud Service have multi-factor authentication activated. This is beneficial to ensure security protocols are being followed within your organization.

```sql+postgres
select
  name,
  id,
  email,
  time_created,
  lifecycle_state,
  is_mfa_activated
from
  oci_identity_user
where
  user_type = 'IDCS';
```

```sql+sqlite
select
  name,
  id,
  email,
  time_created,
  lifecycle_state,
  is_mfa_activated
from
  oci_identity_user
where
  user_type = 'IDCS';
```

### List users who can log in to console
Explore which users have the ability to log in to the console. This can be useful to identify potential security risks and enforce appropriate user permissions.

```sql+postgres
select
  name,
  user_type
from
  oci_identity_user
where
  can_use_console_password;
```

```sql+sqlite
select
  name,
  user_type
from
  oci_identity_user
where
  can_use_console_password = 1;
```

### Details of identity groups attached to users
Explore which user profiles are linked to specific identity groups. This can help in managing user permissions and understanding the distribution of user roles within your organization.

```sql+postgres
select
  oci_identity_user.name as user_name,
  oci_identity_group.name as group_name,
  user_group ->> 'groupId' as group_id
from
  oci_identity_user,
  jsonb_array_elements(user_groups) as user_group
  inner join oci_identity_group ON (oci_identity_group.id = user_group ->> 'groupId' );
```

```sql+sqlite
select
  oci_identity_user.name as user_name,
  oci_identity_group.name as group_name,
  json_extract(user_group.value, '$.groupId') as group_id
from
  oci_identity_user,
  json_each(user_groups) as user_group
  inner join oci_identity_group ON (oci_identity_group.id = json_extract(user_group.value, '$.groupId'));
```