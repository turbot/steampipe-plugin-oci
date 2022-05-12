# Table: oci_limits_quota

You can use quotas to determine how other users allocate Oracle Content Management resources across compartments in Oracle Cloud Infrastructure. Whenever you create an Oracle Content Management instance, the system ensures that your request is within the bounds of the quota for that compartment.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  time_created
from
  oci_limits_quota;
```

### List statements of a quota

```sql
select
  id,
  name,
  statement
from
  oci_limits_quota,
  jsonb_array_elements(statements) as statement;
```

### List active quotas

```sql
select
  id,
  name,
  lifecycle_state,
  time_created,
  statements
from
  oci_limits_quota
where lifecycle_state = 'ACTIVE';
```