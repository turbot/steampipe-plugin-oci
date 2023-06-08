# Table: oci_adm_knowledge_base

Use the Application Dependency Management API to create knowledge bases and vulnerability audits.

## Examples

### Basic info

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base;
```

### List knowledge bases which are not active

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  lifecycle_state <> 'ACTIVE';
```

### List knowledge bases created in last 30 days

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  time_created >= now() - interval '30' day;
```

### List knowledge bases that have not been updated for more than 90 days

```sql
select
  id,
  display_name,
  compartment_id,
  tenant_id,
  lifecycle_state as state
from
  oci_adm_knowledge_base
where
  time_updated < now() - interval '90' day;
```
