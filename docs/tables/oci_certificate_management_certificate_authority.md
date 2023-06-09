# Table: oci_certificate_management_certificate_authority

The metadata details of the certificate authority (CA). This object does not contain the CA contents. For certificate contents see oci_certificates_certificate_authority_bundle.

## Examples

### Basic info

```sql
select
  id,
  name,
  config_type,
  issuer_certificate_authority_id,
  description,
  time_of_deletion,
  kms_key_id,
  lifecycle_details,
  certificate_authority_rules,
  current_version,
  certificate_revocation_list_details,
  subject,
  signing_algorithm,
  lifecycle_state as state
from
  oci_certificate_management_certificate_authority;
```

### List inactive certificate authorities

```sql
select
  id,
  name,
  lifecycle_state,
  issuer_certificate_authority_id,
  current_version
from
  oci_certificate_management_certificate_authority
where
  lifecycle_state <> 'ACTIVE'
```

### List certificates with SHA256_WITH_RSA signed algorithm

```sql
select
  id,
  name,
  current_version,
  signing_algorithm,
  certificate_authority_rules
from
  oci_certificate_management_certificate_authority
where
  signing_algorithm = 'SHA256_WITH_RSA';
```

### List certificates created in the last 30 days

```sql
select
  name,
  id,
  time_created,
  lifecycle_state,
  config_type,
  issuer_certificate_authority_id,
  kms_key_id
from
  oci_certificate_management_certificate_authority
where
  time_created >= now() - interval '30' day;
```

### Get KMS key details of each certificate

```sql
select
  a.id,
  a.name,
  a.kms_key_id,
  k.name as key_name,
  k.vault_id,
  k.current_key_version
from
  oci_certificate_management_certificate_authority as a,
  oci_kms_key as k
where
  k.id = a.kms_key_id;
```

### Get certificate authority rule details

```sql
select
  name,
  id,
  r ->> 'certificateAuthorityMaxValidityDuration' as certificate_authority_max_validity_duration,
  r ->> 'leafCertificateMaxValidityDuration' as leaf_certificate_max_validity_duration,
  r ->> 'ruleType' as rule_type
from
  oci_certificate_management_certificate_authority,
  jsonb_array_elements(certificate_authority_rules) as r;
```

### List certificates that are valid upto a certain date

```sql
select
  name,
  id,
  validity ->> 'timeOfValidityNotAfter' as time_of_validity_not_after
from
  oci_certificate_management_certificate_authority
where
  time_of_validity_not_after <= '2033-06-04T00:00:00Z';
```