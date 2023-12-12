---
title: "Steampipe Table: oci_nosql_table_metric_write_throttle_count_daily - Query OCI NoSQL Tables using SQL"
description: "Allows users to query daily metrics on write throttle count for OCI NoSQL Tables."
---

# Table: oci_nosql_table_metric_write_throttle_count_daily - Query OCI NoSQL Tables using SQL

Oracle NoSQL Database is a scalable, distributed NoSQL database designed to provide highly reliable, flexible, and available data storage across a configurable set of systems. As a part of Oracle's Big Data platform, it delivers predictable latency, is easy to install, manage, and use, and is highly suited to developers. It supports JSON, SQL, and cloud-native API models for application development.

## Table Usage Guide

The `oci_nosql_table_metric_write_throttle_count_daily` table provides insights into the daily metrics on write throttle count for OCI NoSQL Tables. As a database administrator, you can use this table to monitor the write operations that are being throttled daily. This can help you to manage your database performance, by identifying if there are any bottlenecks or if the write capacity needs to be adjusted.

## Examples

### Basic info
Analyze the daily write throttle count metrics to understand trends or identify potential issues in your Oracle NoSQL database. This query helps to monitor the overall performance and health of your database by tracking the minimum, maximum, average, and total write throttle counts over time.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count_daily
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count_daily
order by
  name,
  timestamp;
```

### Intervals where write throttle count exceeded 100 average
Determine the instances when the daily write throttle count surpasses an average of 100. This can be useful to identify potential bottlenecks in data writing processes and facilitate optimization efforts.

```sql+postgres
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count_daily
where
  average > 100
order by
  name,
  timestamp;
```

```sql+sqlite
select
  name,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  oci_nosql_table_metric_write_throttle_count_daily
where
  average > 100
order by
  name,
  timestamp;
```