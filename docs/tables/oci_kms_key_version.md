# Table: oci_kms_key_version

OCI kms key version provides the Key Version resource in Oracle Cloud Infrastructure Kms service.

## Examples

### Basic info

```sql
select
  id,
  key_name,
  lifecycle_state,
  time_created
from
  oci_kms_key_version;
```

### List latest key version of the keys

```sql
select
  key_name,
  max(time_created) as time_created
from
  oci_kms_key_version
group by
  key_name;
```
