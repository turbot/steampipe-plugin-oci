# Table: oci_certificates_management_ca_bundle

Information and metadata for the certificate authority.

## Examples

### Basic info

```sql
select
    id,
    name,
    description,
    lifecycle_details,
    lifecycle_state as state
from
oci_certificates_management_ca_bundle;
```