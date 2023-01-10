# Table: oci_certificate_ca_bundle

The contents of the certificate authority, properties of the certificate authority (and certificate authority version), and user-provided contextual metadata.

## Examples

### Basic info

```sql
select
    id,
    name,
    ca_bundle_pem
from
    oci_certificate_ca_bundle;
```
### Get all Certificate Authority bundles

```sql
select
    cmca.id,
    cca.name,
    cca.ca_bundle_pem
from
    oci_certificate_ca_bundle cca
        inner join
    oci_certificate_management_ca_bundle cmca
    on cca.ca_bundle_id = cmca.id;
```