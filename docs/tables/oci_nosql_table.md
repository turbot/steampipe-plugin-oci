# Table: oci_nosql_table

A NoSQL database service offering on-demand throughput and storage based provisioning that supports JSON, Table and Key-Value datatypes, all with flexible transaction guarantees.

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

### List databases that are not available

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

### List databases with a data storage size greater than 1024 GB

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
