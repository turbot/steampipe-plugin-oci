# Table: oci_identity_tenancy

The tenancy is an Oracle Cloud Account given to you when you register for Oracle Public Cloud (OCI).

## Examples

### Basic info

```sql
select
  name,
  id,
  retention_period_days,
  description
from
  oci_identity_tenancy;
```

### List tenancies where retention period is set to 365 days

```sql
select
  name,
  id,
  retention_period_days,
  home_region_key
from
  oci_identity_tenancy
where
  retention_period_days = 365;
```