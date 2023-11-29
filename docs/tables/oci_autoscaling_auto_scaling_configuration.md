---
title: "Steampipe Table: oci_autoscaling_auto_scaling_configuration - Query OCI Autoscaling Auto Scaling Configurations using SQL"
description: "Allows users to query OCI Auto Scaling Configurations."
---

# Table: oci_autoscaling_auto_scaling_configuration - Query OCI Autoscaling Auto Scaling Configurations using SQL

Auto Scaling in OCI is a cloud service that helps you maintain application availability and allows you to dynamically adjust capacity in response to changing demand patterns. This service ensures that your applications always have the right amount of compute, memory, and other resources they need to stay responsive, even when demand is unpredictable. Auto Scaling automatically adjusts the number of instances in response to changes in load.

## Table Usage Guide

The `oci_autoscaling_auto_scaling_configuration` table provides insights into the Auto Scaling Configurations within OCI Autoscaling. As a DevOps engineer, you can explore configuration-specific details through this table, including policies, resources, and associated metadata. Use it to uncover information about configurations, such as those with specific policies, the resources associated with each configuration, and the current state of each configuration.

## Examples

### Basic info
Determine the areas in which OCI auto-scaling configurations are enabled to gain insights into resource count and creation time. This can be useful for assessing the efficiency and capacity of your auto-scaling setup.

```sql
select
  display_name,
  id,
  is_enabled,
  time_created,
  cool_down_in_seconds,
  max_resource_count
from
  oci_autoscaling_auto_scaling_configuration;
```


### List enabled autoscaling configurations
Explore the configurations that have autoscaling enabled. This is useful to understand the resources that are set to automatically adjust capacity to maintain steady, predictable performance at the lowest possible cost.

```sql
select
  display_name,
  id,
  is_enabled
from
  oci_autoscaling_auto_scaling_configuration
where
  is_enabled;
```


### Get policy details for each autoscaling configuration
This example allows you to analyze all the policies associated with your autoscaling configurations, giving you insights into their status, types, rules, and capacities. This can be beneficial in managing and optimizing your autoscaling configurations for better resource utilization and cost efficiency.

```sql
select
  display_name as autoscaling_configuration_display_name,
  id,
  p ->> 'displayName' as policy_display_name,
  p ->> 'id' as policy_id,
  p ->> 'isEnabled' as policy_is_enabled,
  p ->> 'policyType' as policy_type,
  p ->> 'rules' as rules,
  p ->> 'capacity' as capacity
from
  oci_autoscaling_auto_scaling_configuration,
  jsonb_array_elements(policies) as p
```


### Get resource details for each autoscaling configuration
Explore the details of each autoscaling configuration to understand the associated resources. This can help in managing and optimizing your cloud resources effectively.

```sql
select
  display_name as autoscaling_configuration_display_name,
  id as autoscaling_configuration_id,
  resource ->> 'id' as resource_id,
  resource ->> 'type' as resource_type
from
  oci_autoscaling_auto_scaling_configuration;
```