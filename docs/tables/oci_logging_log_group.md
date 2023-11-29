---
title: "Steampipe Table: oci_logging_log_group - Query OCI Logging Log Groups using SQL"
description: "Allows users to query OCI Logging Log Groups."
---

# Table: oci_logging_log_group - Query OCI Logging Log Groups using SQL

Oracle Cloud Infrastructure (OCI) Logging Log Group is a service that enables you to manage, store, and search logs of all your system, platform, and custom application logs in a centralized location. The Log Group is a resource that represents a group of log streams. It provides you with the ability to categorize and manage related logs together in a structured way.

## Table Usage Guide

The `oci_logging_log_group` table provides insights into Log Groups within Oracle Cloud Infrastructure (OCI) Logging. As a system administrator, you can explore Log Group-specific details through this table, including configurations, associated metadata, and lifecycle states. Utilize it to manage and organize your logs, such as those related to specific applications or environments, and to facilitate efficient log analysis and troubleshooting.

## Examples

### Basic info
Explore the lifecycle state and creation time of different log groups in Oracle Cloud Infrastructure's logging service. This can help you manage and monitor your resources effectively.

```sql
select
  id as log_group_id,
  display_name,
  lifecycle_state,
  time_created
from
  oci_logging_log_group;
```


### List inactive log groups
Explore which log groups are inactive in your OCI logging service. This can be useful for identifying and managing unused resources.

```sql
select
  id as log_group_id,
  display_name,
  lifecycle_state as state,
  time_created
from
  oci_logging_log_group
where
  lifecycle_state = 'INACTIVE';
```