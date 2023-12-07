---
title: "Steampipe Table: oci_nosql_table_metric_storage_utilization_hourly - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query hourly storage utilization metrics of NoSQL Database Tables in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_nosql_table_metric_storage_utilization_hourly - Query OCI NoSQL Database Tables using SQL

The Oracle Cloud Infrastructure (OCI) NoSQL Database service offers a fully managed NoSQL database for modern application development. It's a non-relational, distributed key-value database service that provides predictable latency for data at any scale. It supports JSON document storage and SQL querying, making it flexible for a wide range of applications and use cases.

## Table Usage Guide

The `oci_nosql_table_metric_storage_utilization_hourly` table provides insights into the hourly storage utilization metrics of NoSQL Database Tables within OCI. As a database administrator, you can use this table to monitor and analyze the storage usage pattern of your NoSQL tables on an hourly basis. This can help you in capacity planning, detecting unusual activity, and optimizing the performance of your NoSQL databases.

## Examples

### Basic info
Analyze the settings to understand the storage utilization trends of NoSQL tables over time. This can help in identifying the tables that are consuming more storage, enabling effective resource management.

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
  oci_nosql_table_metric_storage_utilization_hourly
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
  oci_nosql_table_metric_storage_utilization_hourly
order by
  name,
  timestamp;
```

### Intervals where storage utilization exceeded 20GB average
Determine the intervals when the storage utilization surpassed the average of 20GB. This is useful in identifying periods of high storage usage, which can aid in resource planning and cost management.

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
  oci_nosql_table_metric_storage_utilization_hourly
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
  oci_nosql_table_metric_storage_utilization_hourly
where
  average > 20 
order by
  name,
  timestamp;
```