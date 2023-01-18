# Table: oci_application_migration_migration

Application Migration simplifies the migration of applications from Oracle Cloud Infrastructure Classic to Oracle Cloud Infrastructure. You can use Application Migration API to migrate applications, such as Oracle Java Cloud Service, SOA Cloud Service, and Integration Classic instances, to Oracle Cloud Infrastructure.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  description,
  source_id,
  application_name,
  application_type,
  pre_created_target_database_type,
  is_selective_migration,
  service_config,
  application_config,
  lifecycle_details,
  migration_state,
  lifecycle_state as state 
from
  oci_application_migration_migration;
```
