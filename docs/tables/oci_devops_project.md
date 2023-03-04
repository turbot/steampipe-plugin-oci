# Table: oci_devops_project


See: https://registry.terraform.io/providers/oracle/oci/latest/docs/data-sources/devops_project 

## Examples

### Basic info

```sql
select
  display_name,
  id,
  description, 
  namespace,
  state,
  time_created,
  time_updated
from
  oci_devops_project;
```

