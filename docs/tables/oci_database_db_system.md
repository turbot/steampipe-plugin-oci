# Table: oci_database_db_system

Oracle Cloud Infrastructure offers single-node DB systems on either bare metal or virtual machines, and 2-node RAC DB systems on virtual machines.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system;
```

### List db systems that are not available

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  lifecycle_state <> 'AVAILABLE';
```

### List db systems with data storage size greater than 1024 GB

```sql
select
  id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_database_db_system
where
  data_storage_size_in_gbs > 1024;
```
