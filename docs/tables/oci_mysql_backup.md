# Table: oci_mysql_backup

MySQL Database Service supports the following backup types:

FULL: a backup of all data contained in the DB System.
INCREMENTAL: a backup of only the data which has been added or changed since the last FULL backup.

Backups are run in either of the following ways:

Manual: a backup initiated by an action in the console, or request made through the API. Manual backups can be retained for a minimum of 1 day and a maximum of 365 days. Currently, there is a limit of 100 manual backups per tenancy.
Automatic: scheduled backups which run, without any required interaction, at a time of the user's choosing. Automatic backups are retained for between 1 and 35 days. The default retention value is 7 days. Once defined, it is not possible to edit the retention period of an automatic backup.

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
  oci_mysql_backup;
```

### List manual backups

```sql
select
  display_name,
  id,
  creation_type
from
  oci_mysql_backup
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
  oci_mysql_backup
where
  retention_in_days < 90;
```

### Count of backups per DB system

```sql
select
  db_system_id,
  count(*) as backup_count
from
  oci_mysql_backup
group by
  db_system_id;
```
