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
