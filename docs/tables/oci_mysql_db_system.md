# Table: oci_mysql_db_system

A DB System is a logical container for the MySQL instance. It provides an interface enabling management of tasks such as provisioning, backup and restore, monitoring, and so on. It also provides a read/write endpoint enabling you to connect to the MySQL instance using the standard protocols.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system;
```

### List db systems which are not active

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system
where
  lifecycle_state <> 'ACTIVE';
```

### List db systems where backup is not enabled

```sql
select
  id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_mysql_db_system
where
  lifecycle_state = 'ACTIVE'
  and backup_policy -> 'isEnabled' <> 'true';
```
