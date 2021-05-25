# Table: oci_file_storage_file_system

Oracle Cloud Infrastructure File Storage service provides a durable, scalable, secure, enterprise-grade network file system. You can connect to a File Storage service file system from any bare metal, virtual machine, or container instance in your Virtual Cloud Network (VCN).

## Examples

### Basic info

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  availability_domain,
  metered_bytes,
  time_created
from
  oci_file_storage_file_system;
```


## List file systems that are not active

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_file_system
where
  lifecycle_state <> 'ACTIVE';
```


## List cloned file systems

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_file_storage_file_system
where
  is_clone_parent;
```


## List file systems with Oracle managed encryption (default encryption uses Oracle managed encryption keys)

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  time_created
from
  oci_file_storage_file_system
where
  length(kms_key_id) = 0;
```


### List file systems with customer managed encryption keys

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  kms_key_id,
  time_created
from
  oci_file_storage_file_system
where
  length(kms_key_id) <> 0;
```
