# Table: oci_identity_availability_domain

Oracle Cloud Infrastructure is hosted in regions and availability domains. A region is a localized geographic area, and an availability domain is one or more data centers located within a region.

**It will give availability domain details for the subscribed regions only.**

## Examples

### Basic info

```sql
select
  name,
  id
from
  oci_identity_availability_domain;
```
