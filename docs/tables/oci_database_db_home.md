# Table: oci_database_db_home

A directory where Oracle Database software is installed. A bare metal or Exadata DB system can have multiple Database Homes and each Database Home can run a different supported version of Oracle Database. A virtual machine DB system can have only one Database Home.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_home;
```

### List db homes that are not available

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_home
where
  lifecycle_state <> 'AVAILABLE';
```
