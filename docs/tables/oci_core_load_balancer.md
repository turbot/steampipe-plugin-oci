---
title: "Steampipe Table: oci_core_load_balancer - Query OCI Core Load Balancers using SQL"
description: "Allows users to query OCI Core Load Balancers."
---

# Table: oci_core_load_balancer - Query OCI Core Load Balancers using SQL

The Oracle Cloud Infrastructure (OCI) Core Load Balancer is a fully managed service that helps you distribute application traffic across multiple instances in your Virtual Cloud Networks (VCN). It provides automatic scaling and high availability for applications, and supports both internet-facing and private load balancers. The OCI Core Load Balancer allows you to create, manage, and scale load balancers to distribute traffic evenly across your applications.

## Table Usage Guide

The `oci_core_load_balancer` table provides insights into the load balancers within Oracle Cloud Infrastructure (OCI) Core services. As a network engineer or system administrator, you can explore load balancer-specific details through this table, including the configuration, status, and associated instances. Use it to uncover information about load balancers, such as their current operational status, the backend sets, and the associated security groups.

## Examples

### Basic info
Explore which load balancers have been created in your Oracle Cloud Infrastructure, when they were established, and their current operational state. This helps in tracking the lifecycle of your resources and understanding their configurations.

```sql+postgres
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  shape_name
from
  oci_core_load_balancer;
```

```sql+sqlite
select
  display_name,
  id,
  time_created,
  lifecycle_state as state,
  shape_name
from
  oci_core_load_balancer;
```

### List load balancers assigns with public IP address
Explore which load balancers are assigned with a public IP address. This is useful for assessing potential security risks and managing network access.

```sql+postgres
select
  display_name,
  id,
  is_private
from
  oci_core_load_balancer
where
  not is_private;
```

```sql+sqlite
select
  display_name,
  id,
  is_private
from
  oci_core_load_balancer
where
  not is_private;
```