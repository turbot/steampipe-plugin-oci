# Table: oci_certificate_management_certificate_authority_version

The metadata details of the certificate authority (CA) version.

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
    oci_certificate_management_certificate_authority_version;
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
    oci_certificate_management_certificate_authority_version cmcav
        inner join
    oci_certificate_management_certificate_authority cmca
    on cmca.id = cmcav.certificate_authority_id;
```
