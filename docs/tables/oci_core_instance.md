# Table: oci_core_instance

An core instance is a virtual server in the Oracle Cloud Infrastructure cloud.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  shape,
  region
from
  oci_core_instance;
```

### List instances along with the compartment details

```sql
select
  inst.display_name,
  inst.id,
  inst.shape,
  inst.region,
  comp.name as compartment_name
from
  oci_core_instance inst
  inner join
    oci_identity_compartment comp
    on (inst.compartment_id = comp.id)
order by
  comp.name,
  inst.region,
  inst.shape;
```

### Count the number of instances by shape

```sql
select
  shape,
  count(shape) as count
from
  oci_core_instance
group by
  shape;
```

### List instances with shape configuration details

```sql
select
  display_name,
  shape_config_max_vnic_attachments,
  shape_config_memory_in_gbs,
  shape_config_networking_bandwidth_in_gbps,
  shape_config_ocpus,
  shape_config_baseline_ocpu_utilization,
  shape_config_gpus,
  shape_config_local_disks,
  shape_config_local_disks_total_size_in_gbs
from
  oci_core_instance;
```
