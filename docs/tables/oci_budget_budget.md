# Table: oci_budget_budget

A budget can be used to set soft limits on your Oracle Cloud Infrastructure spending. You can set alerts on your budget to let you know when you might exceed your budget, and you can view all of your budgets and spending from one single place in the Oracle Cloud Infrastructure console.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  amount,
  actual_spend,
  lifecycle_state as state
from
  oci_budget_budget;
```

### List budgets that are not active

```sql
select
  display_name,
  id,
  amount,
  actual_spend,
  lifecycle_state as state
from
  oci_budget_budget
where
  lifecycle_state <> 'ACTIVE';
```

### List budgets with actual spend more than 100 percent

```sql
select
  display_name,
  id,
  amount,
  actual_spend,
  lifecycle_state as state
from
  oci_budget_budget
where
  actual_spend > amount;
```
