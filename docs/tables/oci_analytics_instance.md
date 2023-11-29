---
title: "Steampipe Table: oci_analytics_instance - Query OCI Analytics Instances using SQL"
description: "Allows users to query information about OCI Analytics Instances."
---

# Table: oci_analytics_instance - Query OCI Analytics Instances using SQL

Oracle Cloud Infrastructure (OCI) Analytics Instances are part of the OCI Analytics service, which provides an integrated and robust data analytics solution. This service enables users to gather, explore, and analyze data from various sources, transforming it into insightful and actionable business information. It is designed to support all users, from data engineers to business users, with self-service data visualization, enterprise reporting, and augmented analytics.

## Table Usage Guide

The `oci_analytics_instance` table provides insights into OCI Analytics Instances. As a data analyst or business user, you can explore instance-specific details through this table, including current state, capacity, and associated metadata. Utilize it to uncover information about instances, such as those with high capacity, the network configurations of instances, and the verification of instance lifecycle states.

## Examples

### Basic info
Explore which analytics instances are currently active and when they were last modified. This can help in managing resources and identifying instances that may be outdated or unused.

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created,
  time_updated
from
  oci_analytics_instance;
```

### List analytics instances which are not active
Explore which analytics instances are currently inactive. This can be useful for identifying unused resources and potential cost savings.

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  lifecycle_state <> 'ACTIVE';
```

### List analytics instances older than 90 days
Determine the areas in which analytics instances have been active for more than 90 days. This can be useful for assessing system usage and identifying potential areas for optimization or maintenance.

```sql
select
  id,
  name,
  lifecycle_state,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

### List analytics instances with private network endpoint type
Determine the areas in which analytics instances are configured with a private network endpoint. This can be useful for identifying instances that may have heightened security requirements or specific network configurations.

```sql
select
  id,
  name,
  lifecycle_state,
  network_endpoint_details,
  time_created,
  time_updated
from
  oci_analytics_instance
where
  network_endpoint_details ->> 'networkEndpointType' = 'PRIVATE';
```