---
title: "Steampipe Table: oci_identity_api_key - Query OCI Identity API Keys using SQL"
description: "Allows users to query OCI Identity API Keys."
---

# Table: oci_identity_api_key - Query OCI Identity API Keys using SQL

Oracle Cloud Infrastructure's Identity and Access Management (IAM) service lets you control who has access to your cloud resources. You can control what type of access a group of users have and to which specific resources. This is achieved through the use of API Keys, which are used for authentication when making REST API calls to OCI services.

## Table Usage Guide

The `oci_identity_api_key` table provides insights into API Keys within OCI Identity and Access Management (IAM). As a security analyst, explore key-specific details through this table, including its creation time, fingerprint, and associated user ID. Utilize it to uncover information about keys, such as those that are inactive, the users associated with each key, and the verification of key status.

## Examples

### Basic info
Explore which API keys have been created, by whom, and when, to gain insights into your organization's usage and security practices. This can help identify instances where keys may be outdated or potentially misused.

```sql+postgres
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key;
```

```sql+sqlite
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key;
```

### List inactive API keys
Explore which API keys are inactive to ensure your system's security by identifying any unused or potentially compromised keys. This will help maintain the integrity of your system by preventing unauthorized access.

```sql+postgres
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key;
where
  lifecycle_state = 'INACTIVE';
```

```sql+sqlite
select
  key_id,
  key_value,
  user_id,
  user_name,
  time_created,
  fingerprint
from
  oci_identity_api_key
where
  lifecycle_state = 'INACTIVE';
```

### Count API keys by user
Gain insights into how many API keys each user possesses, which can help monitor user access and identify potential security risks.

```sql+postgres
select
  user_id,
  count (*) as api_key_count
from
  oci_identity_api_key
group by
  user_id;
```

```sql+sqlite
select
  user_id,
  count (*) as api_key_count
from
  oci_identity_api_key
group by
  user_id;
```