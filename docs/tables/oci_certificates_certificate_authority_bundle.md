# Table: oci_certificates_certificate_authority_bundle

The contents of the certificate authority, properties of the certificate authority (and certificate authority version), and user-provided contextual metadata.

## Examples

### Basic info

```sql
select
    certificate_authority_id,
    certificate_authority_name,
    serial_number,
    certificate_pem,
    version_number,
    validity,
    stages,
    cert_chain_pem,
    version_name,
    revocation_status
from
oci_certificates_certificate_authority_bundle;
```

### Get all certificate authority bundles
```sql
select
    ccab.certificate_authority_id,
    ccab.certificate_authority_name,
    ccab.serial_number,
    ccab.certificate_pem,
    ccab.version_number,
    ccab.validity,
    ccab.stages,
    ccab.cert_chain_pem,
    ccab.version_name,
    ccab.revocation_status
from
    oci_certificates_certificate_authority_bundle ccab
inner join oci_certificates_management_certificate_authority cmca
on ccab.certificate_authority_id = cmca.id
```