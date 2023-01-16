# Table: oci_core_instance_configuration

An instance configuration is a template that defines the settings to use when creating compute instances.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  time_created,
  region
from
  oci_core_instance_configuration;
```

### List instance configurations along with the compartment details

```sql
select
  inst.display_name,
  inst.id,
  inst.region,
  comp.name as compartment_name
from
  oci_core_instance_configuration inst
  inner join
    oci_identity_compartment comp
    on (inst.compartment_id = comp.id)
order by
  comp.name,
  inst.region;
```

### List instance configurations that are created in last 30 days

```sql
select
  display_name,
  id,
  time_created,
  region
from
  oci_core_instance_configuration
where
  time_created >= now() - interval '30' day;
```
