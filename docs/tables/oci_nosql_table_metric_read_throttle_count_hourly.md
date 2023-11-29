---
title: "Steampipe Table: oci_nosql_table_metric_read_throttle_count_hourly - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query hourly metrics for read throttle count of NoSQL Database Tables in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_nosql_table_metric_read_throttle_count_hourly - Query OCI NoSQL Database Tables using SQL

The Oracle NoSQL Database is a distributed key-value database designed to provide highly reliable, scalable, and available data storage across a configurable set of systems that function as storage nodes. It delivers predictable latency, is easy to use, and offers low total cost of ownership. It enables users to manage and access data across multiple nodes.

## Table Usage Guide

The `oci_nosql_table_metric_read_throttle_count_hourly` table provides insights into the hourly read throttle count metrics of NoSQL Database Tables within Oracle Cloud Infrastructure (OCI). As a database administrator, you can explore table-specific details through this table, including read throttle counts and associated timestamps. Utilize it to monitor and manage the performance of your NoSQL database tables, ensuring optimal usage and avoiding potential bottlenecks.

## Examples

### Basic info
Gain insights into the hourly metrics of NoSQL table read throttle counts, including the minimum, maximum, average, and total counts, to understand the performance and usage patterns. This can be used to optimize resource allocation and improve application performance.

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
  oci_nosql_table_metric_read_throttle_count_hourly
order by
  name,
  timestamp;
```

### Intervals where read throttle count exceeded 100 average
Explore instances where the count of read throttles exceeded an average of 100 within a given hour. This can help in identifying potential bottlenecks in data retrieval and read operations, allowing for better performance optimization.

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
  oci_nosql_table_metric_read_throttle_count_hourly
where
  average > 100
order by
  name,
  timestamp;
```