# Table: oci_certificate_management_association

OCI Certificate Management Association is a feature provided by Oracle Cloud Infrastructure (OCI) that allows you to associate a certificate with a specific resource. In OCI, a certificate represents an SSL/TLS certificate used for securing communication between clients and servers.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  certificates_resource_id,
  association_type,
  time_created
from
  oci_certificate_management_association;
```

### Count numbers of associations by type

```sql
select
  association_type,
  count(id) as numbers_of_association
from
  oci_certificate_management_association
group by
  association_type;
```

### List associations created in the last 10 days

```sql
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificate_management_association
where
  time_created >= now() - interval '10' day;
```

### List associations that are not active

```sql
select
  name,
  id,
  lifecycle_state,
  time_created,
  associated_resource_id
from
  oci_certificate_management_association
where
  lifecycle_state <> 'ACTIVE';
```