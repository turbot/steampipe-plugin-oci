---
title: "Steampipe Table: oci_budget_budget - Query OCI Budgets using SQL"
description: "Allows users to query budget data within Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_budget_budget - Query OCI Budgets using SQL

Oracle Cloud Infrastructure's Budgets service allows users to set thresholds that trigger alerts when cost and usage go beyond the set amount. This service helps users to control costs in their tenancy by providing a way to monitor and manage the consumption of resources. It provides a mechanism to keep track of your OCI spending and ensure it stays within the budget.

## Table Usage Guide

The `oci_budget_budget` table provides insights into budget details within Oracle Cloud Infrastructure's Budgets service. As a finance or operations professional, explore budget-specific details through this table, including the amount, actual spend, and forecasted spend. Utilize it to monitor and manage the consumption of OCI resources, ensuring that spending stays within the defined budget.

## Examples

### Basic info
Gain insights into your budgeting by analyzing the state, amount, and actual spending of your Oracle Cloud Infrastructure resources. This helps in understanding the financial aspects of your resources and planning future expenses accordingly.

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
Discover the segments that are not currently active in your budget. This can be useful to identify areas where resources are being allocated but are not currently in use, helping to optimize your budget distribution.

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
Discover segments where the actual expenditure has exceeded the allocated budget. This is particularly useful for identifying and analyzing areas of overspending to facilitate more effective budget management.

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