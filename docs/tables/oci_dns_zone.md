# Table: oci_dns_zone

A DNS zone holds the trusted DNS records that will reside on Oracle Cloud Infrastructure’s nameservers.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_dns_zone;
```
