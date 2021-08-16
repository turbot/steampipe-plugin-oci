# Table: oci_budget_alert_rule

Budget alert rules can contain rules based on a percentage of your budget or an absolute amount, and on your actual spending or your forecast spending.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  budget_id,
  threshold,
  lifecycle_state as state
from
  oci_budget_alert_rule;
```

### List alert rules that are not active

```sql
select
  display_name,
  id,
  budget_id,
  threshold,
  lifecycle_state as state
from
  oci_budget_alert_rule
where
  lifecycle_state <> 'ACTIVE';
```

### List alert rules with a threshold more than 100 percentage

```sql
select
  display_name,
  id,
  budget_id,
  threshold,
  lifecycle_state as state
from
  oci_budget_alert_rule
where
  threshold > 100
  and threshold_type = 'PERCENTAGE';
```
