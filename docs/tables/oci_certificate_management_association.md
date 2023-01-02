# Table: oci_certificate_management_association

Information about certificates and their associated OCI resources such as Load Balancers and API Gateways.

## Examples

### Basic info

```sql
select
    id,
    name,
    certificates_resource_id,
    associated_resource_id,
    association_type,
    lifecycle_state as state
from
oci_certificate_management_association;
```