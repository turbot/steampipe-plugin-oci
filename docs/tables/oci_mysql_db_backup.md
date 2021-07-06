# Table: oci_mysql_db_backup

Backups help to restore lost data to Cloud SQL instance. Backups protect data from loss or damage.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  backup_type,
  mysql_version
from
  oci_mysql_db_backup;
```


### List manual backups

```sql
select
  display_name,
  id,
  creation_type
from
  oci_mysql_db_backup
where
  creation_type = 'MANUAL';
```


### List backups with retention days less than 90 days

```sql
select
  display_name,
  id,
  retention_in_days
from
  oci_mysql_db_backup
where
  retention_in_days < 90;
```


### Count of backups per db system

```sql
select
  db_system_id,
  count(*) as backup_count
from
  oci_mysql_db_backup
  group by db_system_id;
```