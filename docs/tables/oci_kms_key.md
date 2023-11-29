---
title: "Steampipe Table: oci_kms_key - Query OCI Key Management Service Keys using SQL"
description: "Allows users to query OCI Key Management Service Keys."
---

# Table: oci_kms_key - Query OCI Key Management Service Keys using SQL

Oracle Cloud Infrastructure's Key Management service enables you to manage the cryptographic keys used to protect your data. This service provides centralized key management, key lifecycle management, and cryptographic operations. It allows you to create, import, use, rotate, disable, and delete cryptographic keys.

## Table Usage Guide

The `oci_kms_key` table provides insights into the keys within OCI Key Management Service. As a security engineer, explore key-specific details through this table, including key lifecycle states, creation time, and associated metadata. Utilize it to uncover information about keys, such as those nearing their expiration, the cryptographic algorithm used, and the verification of key usage.

## Examples

### Basic info
Explore the lifecycle state and creation time of your keys in Oracle Cloud Infrastructure's Key Management service. This can help you manage your keys effectively by identifying any keys that are outdated or in an undesirable state.

```sql
select
  id,
  name,
  lifecycle_state,
  time_created,
  vault_name
from
  oci_kms_key;
```

### List keys that are not enabled
Discover the segments that consist of keys not currently enabled. This is useful to identify potential security risks or areas for system optimization.

```sql
select
  id,
  name,
  lifecycle_state,
  vault_name
from
  oci_kms_key
where
  lifecycle_state <> 'ENABLED';
```

### List keys older than 365 days
Determine the areas in which encryption keys have been in use for over a year. This could be useful for identifying outdated security measures and ensuring a regular update cycle for enhanced data protection.

```sql
select
  id,
  name,
  lifecycle_state,
  vault_name
from
  oci_kms_key
where
  time_created <= (current_date - interval '365' day)
order by
  time_created;
```