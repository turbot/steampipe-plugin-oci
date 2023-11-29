---
title: "Steampipe Table: oci_nosql_table_metric_read_throttle_count_daily - Query OCI NoSQL Tables using SQL"
description: "Allows users to query daily metrics for read throttle counts of OCI NoSQL Tables."
---

# Table: oci_nosql_table_metric_read_throttle_count_daily - Query OCI NoSQL Tables using SQL

Oracle NoSQL Database is a fully managed NoSQL database cloud service that provides on-demand throughput and storage-based provisioning. It allows for fast, flexible, and cost-effective data storing and querying, particularly for applications requiring single-digit millisecond latency. It supports both SQL-like querying and JSON documents.

## Table Usage Guide

The `oci_nosql_table_metric_read_throttle_count_daily` table provides insights into the daily read throttle counts of Oracle NoSQL Tables. As a database administrator or developer, you can use this table to understand the read throttle metrics of your NoSQL tables, which can help in performance tuning and cost management. The table can be very useful for identifying trends and patterns in read operations, and for making informed decisions about resource provisioning.

## Examples

### Basic info
Explore which NoSQL tables in your OCI environment have experienced read throttle events. This can help you identify potential performance bottlenecks and optimize your database management.

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_read_throttle_count_daily
order by
  name,
  timestamp;
```

### Intervals where read throttle count exceeded 100 average
Analyze the intervals where the read throttle count exceeded an average of 100. This is useful in identifying potential bottlenecks in your NoSQL database operations, allowing for timely intervention and optimization.

```sql
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_read_throttle_count_daily
where
  average > 100
order by
  name,
  timestamp;
```