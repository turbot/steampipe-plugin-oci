# Table: oci_identity_compartment

Compartments are the primary building blocks you use to organize your cloud resources. You use compartments to organize and isolate your resources to make it easier to manage and secure access to them. For more information about compartments policy, go to [Understanding Compartments](https://docs.oracle.com/en-us/iaas/Content/GSG/Concepts/settinguptenancy.htm#Understa)

## Examples

### Basic info

```sql
select
  name,
  lifecycle_state,
  id,
  compartment_id as parent_compartment
from
  oci.oci_identity_compartment
```

### Details joined with parent compartment

```sql
select
  child.name as name,
  coalesce(parent.name, 'root') as parent_compartment,
  child.id as id,
  coalesce(parent.id, child.tenant_id) as parent_compartment_id
from
  oci_identity_compartment child
  left join oci_identity_compartment parent on (child.compartment_id = parent.id)
order by
  parent.name;
```

### Full path of the compartments

```sql
with recursive compartments as (
select
  name,
  id,
  compartment_id,
  tenant_id,
  name as path,
  name as lastname,
  id as lastid
from 
  oci_identity_compartment
where
  lifecycle_state = 'ACTIVE'

union all

select
  oci_identity_compartment.name,
  oci_identity_compartment.id,
  oci_identity_compartment.compartment_id,
  oci_identity_compartment.tenant_id,
  oci_identity_compartment.name || '\' || compartments.path,
  compartments.lastname,
  compartments.lastid
from 
  oci_identity_compartment join 
  compartments on oci_identity_compartment.id = compartments.compartment_id
)
select
  lastid as compartment_id,
  lastname as name,
  path
from 
  compartments
where 
  compartment_id = tenant_id
order by 
  path;
