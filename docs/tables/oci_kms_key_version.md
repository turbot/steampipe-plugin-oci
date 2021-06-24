# Table: oci_kms_key_version

OCI kms key version provides the Key Version resource in Oracle Cloud Infrastructure Kms service.

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

### List latest key version of the keys

```sql
select
  k.name as key_name,
  max(v.time_created) as latest_key_version_created
from
  oci_kms_key k,
  oci_kms_key_version v
where
  v.key_id = k.id
  and v.management_endpoint = k.management_endpoint
  and v.region = k.region
group by
  key_name;
```
