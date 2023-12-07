---
title: "Steampipe Table: oci_cloud_guard_target - Query OCI Cloud Guard Targets using SQL"
description: "Allows users to query data on Cloud Guard Targets in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_cloud_guard_target - Query OCI Cloud Guard Targets using SQL

Cloud Guard is a security service provided by Oracle Cloud Infrastructure (OCI). It operates as a log and events aggregator that continuously collects, records, and analyzes data, providing a unified view of the security posture of a tenancy. Cloud Guard uses this data to identify any security weak points and actively rectify problems.

## Table Usage Guide

The `oci_cloud_guard_target` table provides insights into Cloud Guard Targets within Oracle Cloud Infrastructure (OCI). As a security or compliance officer, you can utilize this table to explore detailed information about each target, including its status, risk level, and associated resources. This can be particularly useful for monitoring and improving the security posture of your OCI environment, as well as for compliance reporting.

## Examples

### Basic info
Explore the creation times and lifecycle states of your cloud guard targets to better understand their current status and longevity. This can be useful in managing resources and identifying any targets that may require attention or updates.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_target;
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_target;
```

### List targets which are not active
Explore which Cloud Guard targets are not currently active. This is useful for identifying potential security risks or areas that may need maintenance or updates.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state as state
from
  oci_cloud_guard_target
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
  oci_cloud_guard_target
where
  lifecycle_state <> 'ACTIVE';
```