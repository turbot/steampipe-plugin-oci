---
title: "Steampipe Table: oci_nosql_table_metric_read_throttle_count - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query OCI NoSQL Database Tables read throttle count metrics."
---

# Table: oci_nosql_table_metric_read_throttle_count - Query OCI NoSQL Database Tables using SQL

The Oracle Cloud Infrastructure NoSQL Database is a fully managed NoSQL database cloud service that offers developers a high performance, reliable, and cost-effective solution for web-scale data management. The service is easy to use and provides developers with quick and efficient access to their data. It is designed to handle different data models and can be used for many different types of applications, from small, simple applications to large, complex ones.

## Table Usage Guide

The `oci_nosql_table_metric_read_throttle_count` table provides insights into the read throttle count metrics of OCI NoSQL Database Tables. As a database administrator, you can leverage this table to monitor and manage the read throttling on your NoSQL database tables. This can be particularly useful in optimizing the performance and cost efficiency of your NoSQL databases.

## Examples

### Basic info
Explore which NoSQL tables in your OCI environment have experienced read throttling. This can help you identify potential performance issues and optimize your database operations for better efficiency.

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
  oci_nosql_table_metric_read_throttle_count
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
  oci_nosql_table_metric_read_throttle_count
order by
  name,
  timestamp;
```
```

### Intervals where read throttle count exceeded 100 average
Explore instances where the read throttle count exceeded an average of 100, allowing you to identify potential performance issues and optimize accordingly. This can be particularly useful in managing and improving the efficiency of your NoSQL database operations.

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
  oci_nosql_table_metric_read_throttle_count
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
  oci_nosql_table_metric_read_throttle_count
where
  average > 100
order by
  name,
  timestamp;
```