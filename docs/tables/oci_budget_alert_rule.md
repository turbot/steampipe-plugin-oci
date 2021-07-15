# Table: oci_budget_alert_rule

You can set email alerts on your budgets. You can set alerts that are based on a percentage of your budget or an absolute amount, and on your actual spending or your forecast spending.

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

### List alert rules which are not active

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

### List alert rules where threshold more than 100 percentage

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
