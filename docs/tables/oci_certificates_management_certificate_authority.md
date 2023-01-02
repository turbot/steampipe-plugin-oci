# Table: oci_certificates_management_certificate_authority

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
oci_certificates_management_certificate_authority;
```