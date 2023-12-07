---
title: "Steampipe Table: oci_autoscaling_auto_scaling_policy - Query OCI Autoscaling Auto Scaling Policies using SQL"
description: "Allows users to query OCI Autoscaling Auto Scaling Policies."
---

# Table: oci_autoscaling_auto_scaling_policy - Query OCI Autoscaling Auto Scaling Policies using SQL

An OCI Autoscaling Auto Scaling Policy is a feature of Oracle Cloud Infrastructure that automatically adjusts the number of instances in an instance pool based on performance metrics. It allows you to maintain the performance of your applications and manage costs by adjusting the number of instances in response to changes in workload patterns. OCI Autoscaling Auto Scaling Policy helps you ensure that you have the appropriate number of Oracle Cloud Infrastructure Compute instances available to handle the load for your application.

## Table Usage Guide

The `oci_autoscaling_auto_scaling_policy` table provides insights into auto scaling policies within Oracle Cloud Infrastructure Autoscaling. As a Cloud Administrator, explore policy-specific details through this table, including the capacity, policy type, and associated metadata. Utilize it to uncover information about policies, such as their current status, the adjustments made in response to workload changes, and the verification of policy configurations.

## Examples

### Basic info
Assess the elements within your autoscaling policies to understand their current status. This can help in identifying if they are enabled and their capacity, which is useful for managing resource allocation and scaling operations.

```sql+postgres
select
  capacity,
  id,
  display_name,
  is_enabled,
  policy_type
from
  oci_autoscaling_auto_scaling_policy;
```

```sql+sqlite
select
  capacity,
  id,
  display_name,
  is_enabled,
  policy_type
from
  oci_autoscaling_auto_scaling_policy;
```

### List enabled policies
Explore which autoscaling policies are currently active. This can be useful for assessing system performance and identifying areas for optimization.

```sql+postgres
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  is_enabled;
```

```sql+sqlite
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  is_enabled;
```

### List threshold policies
Explore which auto-scaling policies in your OCI environment are of the 'threshold' type. This can help you manage and optimize your resource allocation, by understanding which policies are actively scaling resources based on defined thresholds.

```sql+postgres
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  policy_type = 'threshold';
```

```sql+sqlite
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  policy_type = 'threshold';
```

### List policies older than 30 days
Explore which autoscaling policies have been active for more than 30 days. This can help in assessing the efficiency and relevance of these policies in the current context.

```sql+postgres
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  time_created <= now() - interval '30' day;
```

```sql+sqlite
select
  display_name,
  id,
  policy_type,
  time_created,
  is_enabled
from
  oci_autoscaling_auto_scaling_policy
where
  time_created <= datetime('now', '-30 day');
```

### Get capacity details of each policy
Explore the capacity details of each policy to understand the initial, maximum, and minimum capacity settings. This information can be useful for managing and optimizing the performance of your auto-scaling policies.

```sql+postgres
select
  display_name,
  id,
  capacity ->> 'initial' as initial_capacity,
  capacity ->> 'max' as maximum_capacity,
  capacity ->> 'min' as minimum_capacity
from
  oci_autoscaling_auto_scaling_policy;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(capacity, '$.initial') as initial_capacity,
  json_extract(capacity, '$.max') as maximum_capacity,
  json_extract(capacity, '$.min') as minimum_capacity
from
  oci_autoscaling_auto_scaling_policy;
```

### Get autoscaling configuration details of each policy
Discover the autoscaling configuration details of each policy to better understand your resource management and to assess if your current settings are optimized for your needs. This will help in managing resources efficiently and effectively, ensuring optimal performance.

```sql+postgres
select
  p.display_name,
  p.id,
  p.auto_scaling_configuration_id,
  c.cool_down_in_seconds,
  c.max_resource_count,
  c.min_resource_count
from
  oci_autoscaling_auto_scaling_policy as p,
  oci_autoscaling_auto_scaling_configuration as c
where
  p.auto_scaling_configuration_id = c.id;
```

```sql+sqlite
select
  p.display_name,
  p.id,
  p.auto_scaling_configuration_id,
  c.cool_down_in_seconds,
  c.max_resource_count,
  c.min_resource_count
from
  oci_autoscaling_auto_scaling_policy as p,
  oci_autoscaling_auto_scaling_configuration as c
where
  p.auto_scaling_configuration_id = c.id;
```