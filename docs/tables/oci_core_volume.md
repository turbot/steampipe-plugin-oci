# Table: oci_core_volume

The Oracle Cloud Infrastructure Block Volume service lets you dynamically provision and manage block storage volume

## Examples

### Basic info

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume;
```


### List volumes with faulty state

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume
where
  lifecycle_state = 'FAULTY';
```


### List volumes with memory size greater than 1024 Gb

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_volume
where
  size_in_gbs > 1024;
```


### List volumes with oracle managed encryption( kms_key_id null here indicates that data will be encrypted at rest with an encryption key whose lifecycle management is controlled by Oracle.)

```sql
select
  id as volume_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume
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
  oci_core_volume
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
  oci_core_volume
where
  compartment_id = tenant_id;
```