# Table: oci_vault_secret

OCI vault secrets are credentials such as passwords, certificates, SSH keys, or authentication tokens that you use with Oracle Cloud Infrastructure services.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  id,
  name,
  jsonb_pretty(secret_rules) as rules
from
  oci_vault_secret;
```
