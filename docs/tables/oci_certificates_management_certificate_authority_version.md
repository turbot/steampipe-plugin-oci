# Table: oci_certificates_management_certificate_authority_version

In Oracle Cloud Infrastructure (OCI), the Certificate Management service provides a secure and centralized way to manage SSL/TLS certificates. The certificates managed by OCI Certificate Management do not have versions in the traditional sense. Instead, the service allows you to manage the lifecycle of certificates, including creating, importing, renewing, and revoking certificates.

## Examples

### Basic info

```sql
select
  certificate_authority_id,
  version_number,
  stages,
  serial_number,
  issuer_ca_version_number,
  version_name,
  subject_alternative_names,
  time_of_deletion,
  validity,
  revocation_status
from
  oci_certificates_management_certificate_authority_version;
```

### Get all certificate authority versions

```sql
select
  cmcav.certificate_authority_id,
  cmcav.version_number,
  cmcav.stages,
  cmcav.serial_number,
  cmcav.issuer_ca_version_number,
  cmcav.version_name,
  cmcav.subject_alternative_names,
  cmcav.time_of_deletion,
  cmcav.validity,
  cmcav.revocation_status
from
  oci_certificates_management_certificate_authority_version cmcav
  inner join
    oci_certificates_management_certificate_authority cmca
    on cmca.id = cmcav.certificate_authority_id;
```

### Count certificate versions by certificate authority

```sql
select
  certificate_authority_id,
  count(version_number)
from
  oci_certificates_management_certificate_authority_version
group by
  certificate_authority_id;
```

### List certificate versions created in the last 30 days

```sql
select
  certificate_authority_id,
  version_number,
  serial_number,
  issuer_ca_version_number,
  version_name
from
  oci_certificates_management_certificate_authority_version
where
  time_created >= now() - interval '30' day;
```