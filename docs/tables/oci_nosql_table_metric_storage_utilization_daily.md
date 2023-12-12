---
title: "Steampipe Table: oci_nosql_table_metric_storage_utilization_daily - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query daily metrics for storage utilization of NoSQL Database Tables in Oracle Cloud Infrastructure."
---

# Table: oci_nosql_table_metric_storage_utilization_daily - Query OCI NoSQL Database Tables using SQL

Oracle Cloud Infrastructure's NoSQL Database provides a fully managed NoSQL database service for modern application development. It is designed for applications that require seamless and scalable access to their data, and is ideal for real-time big data workloads. It offers on-demand throughput and storage-based provisioning that supports document, columnar, and key-value data models.

## Table Usage Guide

The `oci_nosql_table_metric_storage_utilization_daily` table provides insights into daily storage utilization metrics of NoSQL Database Tables within Oracle Cloud Infrastructure. As a database administrator or developer, you can use this table to monitor and analyze the storage utilization trends of your NoSQL tables over time. This can help you manage your storage resources more effectively, identify potential issues, and optimize the performance of your NoSQL databases.

## Examples

### Basic info
Analyze the daily storage utilization of NoSQL tables to understand trends such as minimum, maximum, and average usage over time. This can help in optimizing storage allocation, thereby improving cost efficiency and system performance.

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
  oci_nosql_table_metric_storage_utilization_daily
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
  oci_nosql_table_metric_storage_utilization_daily
order by
  name,
  timestamp;
```

### Intervals where storage utilization exceeded 20GB average
Uncover the details of specific periods where the storage utilization exceeded an average of 20GB. This assists in managing storage resources effectively by identifying potential overuse situations.

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
  oci_nosql_table_metric_storage_utilization_daily
where
  average > 20 
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
  oci_nosql_table_metric_storage_utilization_daily
where
  average > 20 
order by
  name,
  timestamp;
```