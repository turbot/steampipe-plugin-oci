# Table: oci_certificate_management_certificate

The details of the certificate. This object does not contain the certificate contents. For certificate contents see oci_certificates_certificate_bundle.

## Examples

### Basic info

```sql
select
  id,
  name,
  config_type,
  issuer_certificate_authority_id,
  description,
  certificate_rules,
  time_of_deletion,
  lifecycle_details,
  current_version,
  subject,
  certificate_revocation_list_details,
  key_algorithm,
  signature_algorithm,
  certificate_profile_type,
  lifecycle_state as state
from
  oci_certificate_management_certificate;
```

### List imported certificates

```sql
select
  name,
  id,
  certificate_rules,
  lifecycle_state,
  config_type,
  signature_algorithm
from
  oci_certificate_management_certificate
where
  config_type = 'IMPORTED';
```

### List failed certificates

```sql
select
  name,
  id,
  current_version,
  subject,
  certificate_revocation_list_details,
  key_algorithm
from
  oci_certificate_management_certificate
where
  lifecycle_state = 'FAILED';
```

### List certificates created in last 30 days

```sql
select
  name,
  id,
  current_version,
  key_algorithm,
  time_created,
  time_of_deletion
from
  oci_certificate_management_certificate
where
  time_created >= now() - interval '30' day;
```

### Get the current version details of each certificate

```sql
select
  name,
  id,
  current_version ->> 'certificateId' as certificate_id,
  current_version ->> 'issuerCaVersionNumber' as issuer_ca_version_number,
  current_version ->> 'revocationStatus' as revocation_status,
  current_version ->> 'serialNumber' as serial_number,
  current_version ->> 'stages' as stages
from
  oci_certificate_management_certificate;
```