---
title: "Steampipe Table: oci_kms_vault - Query OCI Key Management Vaults using SQL"
description: "Allows users to query OCI Key Management Vaults."
---

# Table: oci_kms_vault - Query OCI Key Management Vaults using SQL

OCI Key Management service is a fully managed service that provides centralized key management capabilities, enabling you to manage the entire lifecycle of keys and their associated cryptographic operations. It allows you to create, import, use, rotate, disable, and delete symmetric and asymmetric keys. It also provides centralized key management, key lifecycle management, and cryptographic operations.

## Table Usage Guide

The `oci_kms_vault` table provides insights into Key Management Vaults within Oracle Cloud Infrastructure (OCI). As a security administrator, explore vault-specific details through this table, including the vault type, lifecycle state, and associated metadata. Utilize it to uncover information about vaults, such as those in a particular lifecycle state, the cryptographic endpoint for vaults, and the time of vault creation.

## Examples

### Basic info
Explore which key management vaults are currently active and when they were created. This can help in tracking the lifecycle of your vaults and assessing their configuration for better security management.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  crypto_endpoint,
  management_endpoint,
  vault_type
from
  oci_kms_vault;
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  crypto_endpoint,
  management_endpoint,
  vault_type
from
  oci_kms_vault;
```

### List inactive vaults
Analyze the settings to understand which security vaults are not currently active. This is useful for identifying potential areas of resource optimization and ensuring all inactive vaults are intended to be so.

```sql+postgres
select
  id,
  display_name,
  lifecycle_state,
  vault_type
from
  oci_kms_vault
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  lifecycle_state,
  vault_type
from
  oci_kms_vault
where
  lifecycle_state <> 'ACTIVE';
```

### List virtual private vaults
Explore which key management systems are of the 'virtual private' type to understand the security measures taken in your organization. This can help in assessing the elements within your security infrastructure that are designed for exclusive access.

```sql+postgres
select
  id,
  display_name,
  vault_type
from
  oci_kms_vault
where
  vault_type = 'VIRTUAL_PRIVATE';
```

```sql+sqlite
select
  id,
  display_name,
  vault_type
from
  oci_kms_vault
where
  vault_type = 'VIRTUAL_PRIVATE';
```