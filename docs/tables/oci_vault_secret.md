---
title: "Steampipe Table: oci_vault_secret - Query OCI Vault Secrets using SQL"
description: "Allows users to query OCI Vault Secrets."
---

# Table: oci_vault_secret - Query OCI Vault Secrets using SQL

Oracle Cloud Infrastructure (OCI) Vault is a managed service that centrally manages the encryption of your data. The Vault service is integrated with other OCI services making it easier to manage keys and secrets, and to use them to encrypt data. The Vault service provides centralized key management, key lifecycle management, and cryptographically secure secret management.

## Table Usage Guide

The `oci_vault_secret` table provides insights into secrets within the OCI Vault service. As a security administrator, use this table to explore secret-specific details, including lifecycle details, current version, and associated metadata. Utilize it to uncover information about secrets, such as those nearing the end of their lifecycle, secrets with specific compartment ids, and the verification of secret rules.

## Examples

### Basic info
Explore the basic details of your OCI vault secrets to understand their current lifecycle state and associated vault IDs. This information is useful for managing and tracking the usage of your secrets.

```sql+postgres
select
  name,
  id,
  key_id,
  lifecycle_state,
  vault_id
from
  oci_vault_secret;
```

```sql+sqlite
select
  name,
  id,
  key_id,
  lifecycle_state,
  vault_id
from
  oci_vault_secret;
```

### List secrets in pending deletion state
Identify instances where certain secrets are in a pending deletion state. This can be useful in managing and tracking the lifecycle of your secrets, ensuring no critical data is accidentally lost.

```sql+postgres
select
  name,
  id,
  lifecycle_state
from
  oci_vault_secret
where
  lifecycle_state = 'PENDING_DELETION';
```

```sql+sqlite
select
  name,
  id,
  lifecycle_state
from
  oci_vault_secret
where
  lifecycle_state = 'PENDING_DELETION';
```

### List secret rules
Explore which secret rules are in place within your OCI vault. This is useful for understanding the current security measures and identifying any potential areas for improvement.

```sql+postgres
select
  id,
  name,
  jsonb_pretty(secret_rules) as rules
from
  oci_vault_secret;
```

```sql+sqlite
select
  id,
  name,
  secret_rules as rules
from
  oci_vault_secret;
```