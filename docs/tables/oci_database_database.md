# Table: oci_database_database

An Oracle Database on a bare metal or virtual machine DB system.

## Examples

### Basic info

```sql
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_database;
```

### List databases that are not available

```sql
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_database
where
  lifecycle_state <> 'AVAILABLE';
```
