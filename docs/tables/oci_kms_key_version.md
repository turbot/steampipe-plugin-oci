# Table: oci_kms_key_version

A key version provides new cryptographic material for a master encryption key. The key must be in an ENABLED state to be rotated.

**You must specify a Key ID, Management Endpoint and Region** in a `where` clause (`where key_id='' and management_endpoint='' and region=''`).

## Examples

### Basic info

```sql
select
  v.id as key_version_id,
  k.name as key_name,
  v.lifecycle_state,
  v.time_created as time_created
from
  oci_kms_key k,
  oci_kms_key_version v
where
  v.key_id = k.id
  and v.management_endpoint = k.management_endpoint
  and v.region = k.region;
```

### Get latest key version for all active keys

```sql
select
  k.name as key_name,
  k.lifecycle_state,
  max(v.time_created) as latest_key_version_created,
  k.region,
  coalesce(c.name, 'root') as compartment
from
  oci_kms_key k
  left join oci_identity_compartment c on c.id = k.compartment_id,
  oci_kms_key_version v
where
  v.key_id = k.id
  and v.management_endpoint = k.management_endpoint
  and v.region = k.region
  and k.lifecycle_state = 'ENABLED'
group by
  key_name, k.lifecycle_state, k.region, compartment
```
