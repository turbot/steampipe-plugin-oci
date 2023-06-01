# Table: oci_adm_knowledge_base

Use the Application Dependency Management API to create knowledge bases and vulnerability audits.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state as state 
from
  oci_adm_knowledge_base;
```