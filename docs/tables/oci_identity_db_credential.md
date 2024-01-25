---
title: "Steampipe Table: oci_identity_db_credential - Query OCI Identity DB Credential using SQL"
description: "Allows users to query information about users DB credential in OCI Identity."
---

# Table: oci_identity_db_credential - Query OCI Identity DB Credential using SQL

Oracle Cloud Infrastructure's Identity and Access Management (IAM) service lets you control who has access to your cloud resources. Identity DB credentials refer to the authentication details used to access Oracle databases securely within the Oracle Cloud Infrastructure.

## Table Usage Guide

The `oci_identity_db_credential` table provides insights into users within OCI Identity and Access Management (IAM). As a security administrator, explore user-specific DB credential details through this table, including user ID, name, description, and associated metadata. Utilize it to uncover information about DB credential, such as their state, time of creation, time of expiration, and tenant ID.

## Examples

### Basic info
Discover the segments that highlight DB credential details. This allows for better management and oversight of DB credential lifecycle state create time, expire time.

```sql+postgres
select
  id,
  user_id,
  description,
  time_created,
  lifecycle_state,
  time_expires
from
  oci_identity_db_credential;
```

```sql+sqlite
select
  id,
  user_id,
  description,
  time_created,
  lifecycle_state,
  time_expires
from
  oci_identity_db_credential;
```

###  List Oracle Identity Cloud Service(IDCS) users
Explore which users in the Oracle Identity Cloud Service have multi-factor authentication activated. This is beneficial to ensure security protocols are being followed within your organization.

```sql+postgres
select
  c.id as db_credential_id,
  u.user_type,
  u.is_mfa_activated,
  u.email
from
  oci_identity_db_credential as c,
  oci_identity_user as u
where
  c.user_id = u.id
  and user_type = 'IDCS;
```

```sql+sqlite
select
  c.id as db_credential_id,
  u.user_type,
  u.is_mfa_activated,
  u.email
from
  oci_identity_db_credential as c
  join oci_identity_user as u on c.user_id = u.id
where
  u.user_type = 'IDCS';
```

### List credentials that are set to expire within the next 10 days
Useful for system administrators, security teams, and compliance officers in managing and maintaining the health and security of their OCI (Oracle Cloud Infrastructure) environment.

```sql+postgres
select
  id,
  user_id,
  time_created,
  time_expires
from
  oci_identity_db_credential
where
  time_expires >= now() - interval '10' day;
```

```sql+sqlite
select
  id,
  user_id,
  time_created,
  time_expires
from
  oci_identity_db_credential
where
  time_expires >= datetime('now', '-10 days');
```