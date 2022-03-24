# Table: oci_database_pluggable_database

A pluggable database (PDB) is portable collection of schemas, schema objects, and non-schema objects that appears to an Oracle client as a non-container database.

## Examples

### Basic info

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database;
```

### List pluggable databases that are not available

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database
where
  lifecycle_state <> 'AVAILABLE';
```
