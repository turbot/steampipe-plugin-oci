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

### List boot volumes with faulty state

```sql
select
  id as boot_volume_id
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  lifecycle_state = 'FAULTY';
```

### List boot volumes with size greater than 1024 GB

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_boot_volume
where
  size_in_gbs > 1024;
```

### List boot volumes with Oracle managed encryption
Note: Volumes are encrypted by default with Oracled managed encryption key

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is null;
```

### List boot volumes with customer managed encryption

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_boot_volume
where
  kms_key_id is not null;
```

### List boot volumes created in the root compartment

```sql
select
  id as boot_volume_id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_boot_volume
where
  compartment_id = tenant_id;
```
