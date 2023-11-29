---
title: "Steampipe Table: oci_nosql_table_metric_write_throttle_count - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query write throttle counts for OCI NoSQL Database Tables."
---

# Table: oci_nosql_table_metric_write_throttle_count - Query OCI NoSQL Database Tables using SQL

The OCI NoSQL Database is a fully managed NoSQL database service for modern application development. It offers flexible and cost-effective data management for any scale of real-time, high-volume, key-value workloads. The service ensures predictable low latencies for reads and writes, and offers SQL-like query capabilities.

## Table Usage Guide

The `oci_nosql_table_metric_write_throttle_count` table provides insights into the write throttle counts for OCI NoSQL Database Tables. As a Database Administrator or Developer, explore details about the write throttle counts through this table, which can be useful to understand the performance of your NoSQL database tables. Utilize it to uncover information about the frequency of write throttle events, which can help in optimizing the performance and cost of your NoSQL databases.

## Examples

### Basic info
Explore the metrics related to write throttle counts in NoSQL tables to understand usage patterns and identify potential bottlenecks. This can aid in optimizing table configurations for improved performance.

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
  oci_nosql_table_metric_write_throttle_count
order by
  name,
  timestamp;
```

### Intervals where write throttle count exceeded 100 average
Determine the intervals where the write throttle count exceeded an average of 100. This can help identify periods of high demand or potential performance issues.

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
  oci_nosql_table_metric_write_throttle_count
where
  average > 100
order by
  name,
  timestamp;
```