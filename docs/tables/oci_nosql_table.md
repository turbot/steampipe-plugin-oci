---
title: "Steampipe Table: oci_nosql_table - Query OCI NoSQL Database Tables using SQL"
description: "Allows users to query NoSQL Database Tables in Oracle Cloud Infrastructure (OCI)."
---

# Table: oci_nosql_table - Query OCI NoSQL Database Tables using SQL

Oracle NoSQL Database is a fully managed NoSQL database cloud service that provides on-demand throughput and storage-based provisioning. It is designed for today's most demanding applications that require low latency responses, flexible data models, and elastic scaling for dynamic workloads. With Oracle NoSQL Database, you can focus on application development and deployment without worrying about infrastructure provisioning, patching, and scaling.

## Table Usage Guide

The `oci_nosql_table` table provides insights into NoSQL Database Tables within Oracle Cloud Infrastructure (OCI). As a Database Administrator or Developer, you can explore table-specific details through this table, including the schema, storage, and associated metadata. Utilize it to discover information about tables, such as their capacity, the indexes they use, and the time of their last modification.

## Examples

### Basic info
Explore which NoSQL tables are in use, along with their creation times and lifecycle states, to gain insights into resource allocation and usage patterns. This can help optimize database management and resource planning.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table;
```

### List tables that are not active
Explore which NoSQL tables in your Oracle Cloud Infrastructure are not currently active. This can help in managing resources and identifying any tables that may be underutilized or no longer needed.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table
where
  lifecycle_state <> 'ACTIVE';
```

### List tables with disk storage greater than 1024 GB
Explore which tables are using significant disk storage space to manage resources more effectively. This helps in identifying tables that might be consuming more storage than expected, assisting in efficient resource allocation and cost management.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_nosql_table
where
  cast(table_limits -> 'maxStorageInGBs' as INTEGER) > 1024;
```

### Count child tables for parent tables with children
Analyze the settings to understand the distribution of child tables across parent tables. This helps in assessing the complexity of your NoSQL database structure and can guide optimization efforts.

```sql
select
  t2.name as parent,
  count(t1.*) as child_count
from
  oci_nosql_table t1
  join oci_nosql_table t2 on t1.title like t2.title || '.%' and t1.title <> t2.title
group by
  parent;
```

### Count child tables for parent tables with and without children
Determine the number of child tables linked to each parent table in your NoSQL database. This allows you to understand the complexity and structure of your data, helping in its efficient management and navigation.

```sql
select
  t2.name as parent,
  -- To exclude the parent table from being counted as a child, we subtract 1 from the count.
  count(t1.*) - 1 as child_count
from
  oci_nosql_table t1
  join oci_nosql_table t2 on t1.title like t2.title || '%'
group by
  parent;
```