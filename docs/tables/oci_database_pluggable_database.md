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

### List unavailable pluggable databases

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

### List pluggable databases older than 90 days

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_pluggable_database
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

### List unrestricted pluggable databases

```sql
select
  pdb_name,
  id,
  lifecycle_state,
  is_restricted
from
  oci_database_pluggable_database
where
  not is_restricted;
```