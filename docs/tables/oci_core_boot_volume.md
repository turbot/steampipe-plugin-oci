# Table: oci_core_boot_volume

When you launch a virtual machine (VM) or bare metal instance based on a platform image or custom image, a new boot volume for the instance is created in the same compartment. That boot volume is associated with that instance until you terminate the instance.

## Examples

### Basic info

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume;
```

### List volumes with a faulty state

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  lifecycle_state = 'FAULTY';
```

### List volumes with a memory size greater than 1024 GB

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_boot_volume
where
  size_in_gbs > 1024;
```

### List volumes with Oracle managed encryption (volumes are encrypted by default with Oracled managed encryption keys)

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is null;
```

### List volumes with customer managed encryption

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is not null;
```

### List volumes created in the root compartment

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_boot_volume
where
  compartment_id = tenant_id;
```
