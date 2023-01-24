# Table: oci_core_volume_group

The Oracle Cloud Infrastructure Block Volume service provides you with the capability to group together multiple volumes in a volume group. A volume group can include both types of volumes, boot volumes, which are the system disks for your compute instances, and block volumes for your data storage.

## Examples

### Basic info

```sql
select
  id as volume_group_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_group;
```

### List volume groups with a faulty state

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_core_volume_group
where
  lifecycle_state = 'FAULTY';
```

### List volume groups with size greater than 1024 GB

```sql
select
  id,
  display_name,
  lifecycle_state,
  size_in_gbs
from
  oci_core_volume_group
where
  size_in_gbs > 1024;
```

### List volume groups created in the root compartment

```sql
select
  id,
  display_name,
  lifecycle_state,
  tenant_id,
  compartment_id
from
  oci_core_volume_group
where
  compartment_id = tenant_id;
```

### List volume groups that are created in the last 30 days

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created,
  size_in_gbs
from
  oci_core_volume_group
where
  time_created >= now() - interval '30' day;
```

### Get volume details for the volume groups

```sql
select
  g.id as volume_group_id,
  g.display_name as volume_group_diplay_name,
  v.id as volume_id,
  v.auto_tuned_vpus_per_gb
from
  oci_core_volume_group as g,
  oci_core_volume as v,
  jsonb_array_elements_text(volume_ids) as i
where
  v.id = i;
```