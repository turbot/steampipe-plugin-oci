# Table: oci_certificate_management_certificate_version

The details of the certificate version

## Examples

### Basic info

```sql
select
  certificate_id,
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
  oci_certificate_management_certificate_version;
```

### Get all certificate versions
```sql
select
  cmcv.certificate_id,
  cmcv.version_number,
  cmcv.stages,
  cmcv.serial_number,
  cmcv.issuer_ca_version_number,
  cmcv.version_name,
  cmcv.subject_alternative_names,
  cmcv.time_of_deletion,
  cmcv.validity,
  cmcv.revocation_status 
from
  oci_certificate_management_certificate_version cmcv 
  inner join
    oci_certificate_management_certificate cmc 
    on cmcv.certificate_id = cmc.id;
```
