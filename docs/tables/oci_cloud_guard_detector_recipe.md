---
title: "Steampipe Table: oci_cloud_guard_detector_recipe - Query OCI Cloud Guard Detector Recipes using SQL"
description: "Allows users to query OCI Cloud Guard Detector Recipes."
---

# Table: oci_cloud_guard_detector_recipe - Query OCI Cloud Guard Detector Recipes using SQL

OCI Cloud Guard is a security service that identifies potential security weaknesses and consolidates this information into a single dashboard. It provides a comprehensive view of the security and compliance status of your resources in Oracle Cloud Infrastructure. Detector Recipes in Cloud Guard contain the configurations of detectors and the conditions that cause them to trigger.

## Table Usage Guide

The `oci_cloud_guard_detector_recipe` table provides insights into the Detector Recipes within OCI Cloud Guard. As a security analyst, explore recipe-specific details through this table, including the associated managed list types, detector rules, and associated metadata. Utilize it to uncover information about detector recipes, such as their status, the conditions that trigger them, and the actions taken when those conditions are met.

## Examples

### Basic info
Explore which cloud guard detector recipes have been created, their respective IDs, when they were created, and their current lifecycle states. This information can help in managing and tracking the status of your cloud guard detector recipes.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe;
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe;
```

### List detector recipes which are not active
Explore which detector recipes in your Cloud Guard are inactive. This is useful for identifying potential security blind spots in your system.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_detector_recipe
where
  lifecycle_state <> 'ACTIVE';
```

### List detector recipes with password related rules
Discover the segments that have password-related rules within the detector recipes. This query is useful for assessing security measures in place and ensuring rules regarding password age and complexity are being enforced.

```sql+postgres
select
  name,
  e ->> 'detectorRuleId' as Rule_name,
  e -> 'details' ->> 'isEnabled' as status
from
  oci_cloud_guard_detector_recipe,
  jsonb_array_elements(effective_detector_rules) as e
where
  e ->> 'detectorRuleId' = 'PASSWORD_TOO_OLD'
  or e ->> 'detectorRuleId' = 'PASSWORD_POLICY_NOT_COMPLEX';
```

```sql+sqlite
select
  name,
  json_extract(e.value, '$.detectorRuleId') as Rule_name,
  json_extract(json_extract(e.value, '$.details'), '$.isEnabled') as status
from
  oci_cloud_guard_detector_recipe,
  json_each(effective_detector_rules) as e
where
  json_extract(e.value, '$.detectorRuleId') = 'PASSWORD_TOO_OLD'
  or json_extract(e.value, '$.detectorRuleId') = 'PASSWORD_POLICY_NOT_COMPLEX';
```