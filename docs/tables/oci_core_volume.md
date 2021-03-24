# Table: oci_core_volume

The Oracle Cloud Infrastructure Block Volume service lets you dynamically provision and manage block storage volume

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_new.oci_core_volume;
```


### List of volumes which are in faulty state

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_new.oci_core_volume
where
  lifecycle_state = 'FAULTY'
```


### List of volumes with memory size greater tham 1024 Gb

```sql
select
  id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_new.oci_core_volume
where
  size_in_gbs > 1024
```