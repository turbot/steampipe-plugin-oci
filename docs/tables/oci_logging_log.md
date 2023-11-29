---
title: "Steampipe Table: oci_logging_log - Query OCI Logging Logs using SQL"
description: "Allows users to query OCI Logging Logs."
---

# Table: oci_logging_log - Query OCI Logging Logs using SQL

Oracle Cloud Infrastructure's (OCI) Logging service is a highly scalable and fully managed single pane of glass for all the logs in your tenancy. The service helps you to manage and analyze logs from your resources in Oracle Cloud Infrastructure, your applications, and your on-premises resources. This makes it easier to monitor, troubleshoot, and react to operational and security issues.

## Table Usage Guide

The `oci_logging_log` table provides insights into logs within Oracle Cloud Infrastructure's (OCI) Logging service. As a system administrator, this table is useful to explore log-specific details, including the log group it belongs to, the log type, and the configuration details. Utilize it to uncover information about logs, such as those with specific configurations, the relationships between logs and log groups, and the status of each log.

## Examples

### Basic info
Explore which logs have been created within your Oracle Cloud Infrastructure (OCI) environment, and assess their lifecycle state to understand if they are active or deleted. This can be useful for managing and tracking your OCI resources.

```sql
select
  id,
  log_group_id,
  name,
  lifecycle_state,
  time_created
from
  oci_logging_log;
```

### List inactive logs
Identify logs that are currently inactive. This can be useful in managing system resources or troubleshooting system issues by focusing on logs that are not actively recording data.

```sql
select
  id,
  name,
  lifecycle_state as state,
  time_created
from
  oci_logging_log
where
  lifecycle_state = 'INACTIVE';
```

### List VCN subnets with flow logging enabled
Assess the elements within your network by identifying active subnets that have flow logging enabled. This can help enhance network security and troubleshooting by providing visibility into traffic patterns and potential anomalies.

```sql
select
  configuration -> 'source' ->> 'resource' as subnet_id,
  configuration -> 'source' ->> 'service' as service,
  lifecycle_state
from
  oci_logging_log
where
  configuration -> 'source' ->> 'service' = 'flowlogs'
  and lifecycle_state = 'ACTIVE';
```