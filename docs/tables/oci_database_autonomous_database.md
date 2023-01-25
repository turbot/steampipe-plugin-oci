# Table: oci_database_autonomous_database

Oracle Cloud Infrastructure's Autonomous Database is a fully managed, preconfigured database environment with four workload types available:

- Autonomous Transaction Processing
- Autonomous Data Warehouse
- Oracle APEX Application Development
- Autonomous JSON Database

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

### Get KMS key details for the databases

```sql
select
  d.db_name,
  d.display_name,
  d.kms_key_id,
  k.name as key_name,
  k.algorithm as key_algorithm,
  k.current_key_version,
  k.protection_mode
from
  oci_database_autonomous_database as d,
  oci_kms_key as k
where
  k.id = d.kms_key_id;
```

### Get KMS vault details for the databases

```sql
select
  d.db_name,
  d.display_name,
  d.vault_id,
  v.display_name as vault_display_name,
  v.crypto_endpoint,
  v.vault_type,
  v.management_endpoint
from
  oci_database_autonomous_database as d,
  oci_kms_vault as v
where
  v.id = d.vault_id;
```