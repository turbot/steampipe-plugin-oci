# Table: oci_certificate_certificate_bundle_version

The contents of the certificate, properties of the certificate (and certificate version), and user-provided contextual metadata.

## Examples

### Basic info

```sql
select
  certificate_id,
  certificate_name,
  version_number,
  serial_number,
  validity,
  stages,
  certificate_pem,
  cert_chain_pem,
  private_key_pem,
  version_name,
  revocation_status,
  certificate_bundle_type
from
  oci_certificate_certificate_bundle_version;
```

### Get all certificate bundles

```sql
select
  ccb.certificate_id,
  ccb.certificate_name,
  ccb.version_number,
  ccb.serial_number,
  ccb.validity,
  ccb.stages,
  ccb.certificate_pem,
  ccb.cert_chain_pem,
  ccb.private_key_pem,
  ccb.version_name,
  ccb.revocation_status,
  ccb.certificate_bundle_type
from
    oci_certificate_certificate_bundle_version ccb
inner join oci_certificate_management_certificate cmc
on ccb.certificate_id = cmc.id;
```

### Count bundles by certificate

```sql
select
  certificate_id,
  certificate_name,
  count(version_number) as number_of_versions
from
  oci_certificate_certificate_bundle_version
group by
  certificate_id,
  certificate_name;
```