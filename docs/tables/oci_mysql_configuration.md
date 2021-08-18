# Table: oci_mysql_configuration

DB System configurations are collections of variables which define the operation of the MySQL DB System.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration;
```

### List deleted configurations

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_configuration 
where
  lifecycle_state = 'DELETED'
```