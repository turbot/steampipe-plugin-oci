# Table: oci_audit_configuration

By default, Audit logs are retained for 90 days. You can configure log retention for up to 365 days. You can edit the log retention period in the tenancy details page.

## Examples

### Basic info

```sql
select
  retention_period_days
from
  oci_audit_configuration;
```
