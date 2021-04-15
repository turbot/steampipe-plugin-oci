# Table: oci_kms_vault

Oracle Cloud Infrastructure Vault is a managed service that lets you centrally manage the encryption keys that protect your data and the secret credentials that you use to securely access resources.

## Examples

### Basic info

```sql
select
  id ,
  display_name,
  lifecycle_state,
  time_created,
  crypto_endpoint,
  management_endpoint,
  vault_type
from
  oci_kms_vault;
```

### List of vaults where state is not active

```sql
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


### List of vaults where vault type us virtual private

```sql
select
  id,
  display_name,
  vault_type
from
  oci_kms_vault
where
  vault_type = 'VIRTUAL_PRIVATE';
```