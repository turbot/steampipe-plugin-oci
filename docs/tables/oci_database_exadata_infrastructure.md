---
title: "Steampipe Table: oci_database_exadata_infrastructure - Query OCI Database Exadata Infrastructures using SQL"
description: "Allows users to query OCI Database Pluggable Databases."
---

# Table: oci_database_exadata_infrastructure - Query OCI Database Exadata Infrastructures using SQL

A Database Exadata Infrastructure is a specialized Oracle cloud service that provides high-performance and highly available database infrastructure for running Oracle Databases. Exadata Infrastructure is designed to deliver superior database performance and scalability by leveraging Oracle Exadata technology, which is optimized for Oracle Database workloads.

## Table Usage Guide

The `oci_database_exadata_infrastructure` table provides insights into Database Exadata Infrastructure within Oracle Cloud Infrastructure Database service. As a database administrator, you can explore specific details about each Exadata Infrastructure through this table, including its configuration, status, and associated metadata. Utilize it to uncover information about infrastructures, such as their database performance, making them suitable for demanding workloads such as OLTP (Online Transaction Processing), data warehousing, and analytics.

## Examples

### Basic info

Explore which infrastructures are in different stages of their lifecycle and when they were created. This allows you to assess the overall health and status of your infrastructures.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  time_zone,
  cpus_enabled,
  max_cpu_count
from
  oci_database_exadata_infrastructure;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  time_zone,
  cpus_enabled,
  max_cpu_count
from
  oci_database_exadata_infrastructure;
```

### List failed exadata infrastructures

This query can be used to identify any infrastructures in the Oracle Cloud Infrastructure that have failed to initialize. By doing so, it allows users to quickly pinpoint and address any issues that may be disrupting their infrastructure operations.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  shape
from
  oci_database_exadata_infrastructure
where
  lifecycle_state = 'FAILED';
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  time_created,
  shape
from
  oci_database_exadata_infrastructure
where
  lifecycle_state = 'FAILED';
```

### List infrastructures older than 90 days

Discover the segments that consist of exadata infrastructure older than 90 days. This is useful for identifying potential areas for system optimization or data archiving.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_exadata_infrastructure
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  time_created
from
  oci_database_exadata_infrastructure
where
  time_created <= date('now','-90 day')
order by
  time_created;
```

### Top 5 exadata infrastructures with max CPUs

Discover the top 5 Exadata Infrastructures with maximum CPU power, perfect for high-performance workloads.

```sql+postgres
select
  display_name,
  id,
  lifecycle_state,
  shape,
  cpus_enabled,
  max_cpu_count
from
  oci_database_exadata_infrastructure
order by
  max_cpu_count desc
limit
  1;
```

```sql+sqlite
select
  display_name,
  id,
  lifecycle_state,
  shape,
  cpus_enabled,
  max_cpu_count
from
  oci_database_exadata_infrastructure
order by
  max_cpu_count desc
limit
  1;
```
