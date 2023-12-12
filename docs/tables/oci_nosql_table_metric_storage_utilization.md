---
title: "Steampipe Table: oci_nosql_table_metric_storage_utilization - Query OCI NoSQL Tables using SQL"
description: "Allows users to query the storage utilization metrics of OCI NoSQL Tables."
---

# Table: oci_nosql_table_metric_storage_utilization - Query OCI NoSQL Tables using SQL

Oracle NoSQL Database is a scalable, low-latency key-value database for real-time big data workloads. It is designed to provide highly reliable, flexible, and available data storage across a configurable set of systems. This database supports SQL-like query capabilities and transactional consistency, making data easy to manage and analyze.

## Table Usage Guide

The `oci_nosql_table_metric_storage_utilization` table provides insights into storage utilization metrics of Oracle NoSQL tables. As a database administrator, you can leverage this table to monitor and manage the storage usage of your NoSQL tables. It helps you identify tables that are nearing their storage capacity, allowing you to take proactive measures to ensure optimal performance and avoid potential issues.

## Examples

### Basic info
Analyze the storage utilization of your NoSQL tables to understand usage patterns over time. This can help you optimize resource allocation and manage costs more effectively.

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
  oci_nosql_table_metric_storage_utilization
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
  oci_nosql_table_metric_storage_utilization
order by
  name,
  timestamp;
```

### Intervals where storage utilization exceeded 20GB average
Determine the periods where your storage usage exceeded an average of 20GB. This is useful to identify potential resource overutilization and optimize your storage management.

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
  oci_nosql_table_metric_storage_utilization
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
  oci_nosql_table_metric_storage_utilization
where
  average > 20 
order by
  name,
  timestamp;
```