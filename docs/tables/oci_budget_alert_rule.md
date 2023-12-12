---
title: "Steampipe Table: oci_budget_alert_rule - Query OCI Budgets Alert Rules using SQL"
description: "Allows users to query OCI Budgets Alert Rules."
---

# Table: oci_budget_alert_rule - Query OCI Budgets Alert Rules using SQL

Oracle Cloud Infrastructure (OCI) Budgets service allows you to set budget thresholds and send notifications when the thresholds are exceeded. Alert rules can be set up to track spending and usage, allowing users to manage their costs and consumption effectively. Alert rules are associated with a specific budget and trigger notifications based on the thresholds defined.

## Table Usage Guide

The `oci_budget_alert_rule` table provides insights into alert rules within OCI Budgets service. As a financial analyst or cloud cost manager, you can explore details about each alert rule through this table, including the associated budget, threshold, and type of alert. Use this table to track and manage spending, understand cost trends, and ensure budget compliance.

## Examples

### Basic info
Explore the financial boundaries of your project by identifying the alert rules in your budget. This query can help you assess the financial health and lifecycle state of your project.

```sql+postgres
select
  display_name,
  id,
  budget_id,
  threshold,
  lifecycle_state as state
from
  oci_budget_alert_rule;
```

```sql+sqlite
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
Explore which budget alert rules within your OCI environment are not currently active. This is useful for identifying potential savings or areas where budget tracking may not be effectively implemented.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which budget alert rules have been set with a threshold exceeding 100 percent. This can be useful to identify potential errors or overly conservative settings in your financial monitoring.

```sql+postgres
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

```sql+sqlite
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