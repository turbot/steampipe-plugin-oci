# Table: oci_mysql_configuration_custom

DB System configurations are collections of variables which define the operation of the MySQL DB System.

Configurations have a default set of variables assigned to them, user and system variables. To add variables, you must create a new configuration with the desired variable definitions, or copy an existing configuration, edit it accordingly, and edit the DB System to use the new configuration.

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
  oci_mysql_configuration_custom;
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
  oci_mysql_configuration_custom 
where
  lifecycle_state = 'DELETED';
```
