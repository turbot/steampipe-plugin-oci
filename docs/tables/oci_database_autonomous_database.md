# Table: oci_database_autonomous_database

Oracle Cloud Infrastructure's Autonomous Database is a fully managed, preconfigured database environment with four workload types available, which are:

- Autonomous Transaction Processing
- Autonomous Data Warehouse
- Oracle APEX Application Development
- Autonomous JSON Database

You do not need to configure or manage any hardware or install any software.

## Examples

### Basic info

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database;
```

### List databases which are not available

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database
where
  lifecycle_state <> 'AVAILABLE';
```

### List databases where data storage size greater than 1024 GB

```sql
select
  db_name,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_autonomous_database
where
  data_storage_size_in_gbs > 1024;
```
