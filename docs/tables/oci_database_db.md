# Table: oci_database_db

An Oracle Database on a bare metal or virtual machine DB system. All single-node Oracle RAC DB systems support the following Oracle Database editions:

- Standard Edition
- Enterprise Edition
- Enterprise Edition - High Performance
- Enterprise Edition - Extreme Performance

Two-node Oracle RAC DB systems require Oracle Enterprise Edition - Extreme Performance.

## Examples

### Basic info

```sql
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db;
```

### List databases that are not available

```sql
select
  db_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_db
where
  lifecycle_state <> 'AVAILABLE';
```
