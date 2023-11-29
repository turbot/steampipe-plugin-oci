---
title: "Steampipe Table: oci_mysql_db_system_metric_memory_utilization - Query OCI MySQL DB System Metrics using SQL"
description: "Allows users to query MySQL DB System Memory Utilization Metrics."
---

# Table: oci_mysql_db_system_metric_memory_utilization - Query OCI MySQL DB System Metrics using SQL

Oracle Cloud Infrastructure's MySQL Database Service is a fully managed database service that lets developers and database administrators focus on application development and business innovation. It provides a highly available, scalable, and secure MySQL Database environment, with the operational benefits of the cloud. The service is built on MySQL Enterprise Edition and powered by Oracle Cloud Infrastructure.

## Table Usage Guide

The `oci_mysql_db_system_metric_memory_utilization` table provides insights into the memory utilization metrics of MySQL DB Systems in Oracle Cloud Infrastructure (OCI). As a database administrator, you can use this table to monitor and analyze the memory usage patterns of your MySQL databases, which can help you optimize performance and resource allocation. Additionally, it can assist in identifying potential issues related to memory utilization, enabling proactive troubleshooting and maintenance.

## Examples

### Basic info
Explore the performance of your MySQL database system by examining the memory utilization metrics. This allows you to identify periods of high or low usage, aiding in resource planning and optimization.

```sql
select
  id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization
order by
  id,
  timestamp;
```

### Memory Utilization Over 80% average
Explore instances where the average memory utilization exceeds 80% in your MySQL database system. This can help you identify potential performance issues and optimize resource allocation for better system performance.

```sql
select
  id,
  timestamp,
  round(minimum::numeric,2) as min_memory,
  round(maximum::numeric,2) as max_memory,
  round(average::numeric,2) as avg_memory,
  sample_count
from
  oci_mysql_db_system_metric_memory_utilization
where average > 80
order by
  id,
  timestamp;
```