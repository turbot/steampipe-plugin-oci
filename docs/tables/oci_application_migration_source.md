# Table: oci_application_migration_source

The properties that define a source. Source refers to the source environment from which you migrate an application to Oracle Cloud Infrastructure.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source;
```

### List all inactive sources

```sql
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source
where
  lifecycle_state <> 'ACTIVE';
```

### List all sources created in the last 30 days

```sql
select
  id,
  display_name,
  description,
  lifecycle_details,
  lifecycle_state as state,
  source_details ->> 'computeAccount' as source_account_name,
  source_details ->> 'region' as source_account_region,
  source_details ->> 'type' as source_account_type
from
  oci_application_migration_source
where
  time_created >= now() - interval '30' day;
```

### List all migrations under a particular source

```sql
select
  s.id as source_id,
  s.display_name as source_name,
  m.id as migration_id,
  m.display_name as migration_name,
  m.lifecycle_details,
  m.lifecycle_state as state
from
  oci_application_migration_source as s,
  oci_application_migration_migration as m
where
  s.id = m.source_id
  and s.display_name = 'source-1';
```