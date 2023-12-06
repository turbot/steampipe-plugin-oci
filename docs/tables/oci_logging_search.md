---
title: "Steampipe Table: oci_logging_search - Query OCI Logging Search using SQL"
description: "Allows users to query OCI Logging Search data."
---

# Table: oci_logging_search - Query OCI Logging Search using SQL

Oracle Cloud Infrastructure's Logging Search is a fully managed, scalable, and highly available service that allows you to centralize, analyze, and monitor your logs. Logging Search provides a user-friendly interface for you to search, filter, and analyze logs from all your resources in one place. It is designed to help you diagnose, troubleshoot, and monitor your applications and infrastructure.

## Table Usage Guide

The `oci_logging_search` table provides insights into the Logging Search within Oracle Cloud Infrastructure (OCI). As a DevOps engineer, explore search-specific details through this table, including log groups, log content, and associated metadata. Utilize it to uncover information about logs, such as those related to specific resources, the content of the logs, and the verification of log groups.

**Important Notes**
- By default, this table will provide data for the last 24hrs. You can pass the `timestamp` value in the following ways to fetch data in a range.
  - timestamp >= ‘2023-03-11T00:00:00Z’ and timestamp <= ‘2023-03-15T00:00:00Z’
  - timestamp between ‘2023-03-11T00:00:00Z’ and ‘2023-03-15T00:00:00Z’
  - timestamp >= now() - interval '7 days' (The data will be fetched for the last 7 days)
  - timestamp > ‘2023-03-15T00:00:00Z’ (The data will be fetched from the provided time to the current time)
  - timestamp < ‘2023-03-15T00:00:00Z’ (The data will be fetched from one day before the provided time to the provided time)
- We recommend specifying optional quals in the query to optimize the table output. Optional quals are supported for the following columns:
  - `log_group_name`
  - `log_name`
  - `search_query`
  - `timestamp`

## Examples

### Show log entries of the last 24 hrs
Gain insights into the recent activities within your system by viewing log entries from the past 24 hours. This can help you monitor system performance, identify potential issues, and maintain system health.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search;
```

### Show log entries of the last 3 days
Explore the recent activities in your system by analyzing the log entries from the past three days. This can be particularly useful for monitoring system health, detecting anomalies, and troubleshooting issues.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  timestamp >= now() - interval '3 days';
```

### Show log entries that are related to DatabaseService
Explore log entries associated with a specific service to gain insights into system activity and performance. This can help identify potential issues or anomalies related to the service, enhancing overall system monitoring and management.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  log_content_type like '%DatabaseService%';
```

### Show log entries from us-ashburn-1 region
Explore the log entries originating from a specific region to gain insights into system performance and potential issues. This can be useful for troubleshooting or optimizing operations within that region.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  region = 'us-ashburn-1';
```

### Show log entries from a specific log group
Analyze the entries from a particular log group to gain insights into specific events or issues. This is beneficial in troubleshooting and understanding the operational behavior of your system.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  log_group_name = 'test-bucket';
```

### Show log entries from a specific search query
Explore log entries associated with a particular search query. This can help in analyzing patterns, identifying issues, and gaining insights into specific activities within a given timeframe or region.

```sql
select
  log_content_id,
  log_content_source,
  log_content_type,
  timestamp,
  region,
  compartment_id,
  tenant_id,
  jsonb_pretty(log_content) as log_content
from
  oci_logging_search
where
  search_query = 'search "ocid1.tenancy.oc1..aaaaaaaahnm7gleh5soecx3hoz4p4h2q37cyljaq/test" | sort by datetime desc';
```