---
title: "Steampipe Table: oci_nosql_table_metric_write_throttle_count_hourly - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query hourly metrics for write throttle count of OCI NoSQL Database Tables."
---

# Table: oci_nosql_table_metric_write_throttle_count_hourly - Query OCI NoSQL Database Tables using SQL

The Oracle Cloud Infrastructure (OCI) NoSQL Database service is a fully managed NoSQL database cloud service for modern application development. It's designed for applications that require seamless and fast access to their data, and it provides flexible and cost-effective data management for any scale of data. The service is easy to use, elastic, and available for several Oracle Cloud Infrastructure regions worldwide.

## Table Usage Guide

The `oci_nosql_table_metric_write_throttle_count_hourly` table provides insights into the hourly metrics for write throttle count of OCI NoSQL Database Tables. As a database administrator, you can use this table to monitor the frequency of write throttling events, which could indicate potential performance issues or bottlenecks in your NoSQL database tables. This data can be instrumental in optimizing the performance and efficiency of your database operations.

## Examples

### Basic info
Explore the performance metrics of NoSQL databases in Oracle Cloud Infrastructure. This query provides insights into the write throttle count on an hourly basis, which can help in identifying any potential bottlenecks or performance issues.

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
  oci_nosql_table_metric_write_throttle_count_hourly
order by
  name,
  timestamp;
```

### Intervals where write throttle count exceeded 100 average
Analyze the settings to understand intervals where the average write throttle count surpassed 100. This could be beneficial in identifying potential performance issues in your NoSQL databases.

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
  oci_nosql_table_metric_write_throttle_count_hourly
where
  average > 100
order by
  name,
  timestamp;
```