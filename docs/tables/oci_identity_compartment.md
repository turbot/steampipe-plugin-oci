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
  parent.name as parent_compartment,
  child.id as id,
  parent.id as parent_compartment_id
from
  oci_identity_compartment child
  left join
    oci_identity_compartment parent
    on (child.compartment_id = parent.id)
order by
  parent.name;
```
