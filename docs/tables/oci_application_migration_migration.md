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

### List all inactive migrations

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
  oci_application_migration_migration
where
  lifecycle_state <> 'ACTIVE';
```

### List all migrations created in the last 30 days

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
  oci_application_migration_migration
where
  time_created >= now() - interval '30' day;
```

### List source details of a particular migration

```sql
select
  m.id as migration_id,
  m.display_name as migration_name,
  s.id as source_id,
  s.display_name as source_name,
  s.lifecycle_details,
  s.lifecycle_state as state,
  s.source_details ->> 'computeAccount' as source_account_name,
  s.source_details ->> 'region' as source_account_region,
  s.source_details ->> 'type' as source_account_type
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and m.display_name = 'migration-1';
```

### List successfully migrated applications

```sql
select
  id as migration_id,
  display_name as migration_name,
  application_config,
  application_name,
  application_type
from
  oci_application_migration_migration
where
  migration_state = 'MIGRATION_SUCCEEDED';
```

### List application migrations that have selective-migration enabled

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
  oci_application_migration_migration
where
  is_selective_migration;
```