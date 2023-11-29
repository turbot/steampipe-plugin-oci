---
title: "Steampipe Table: oci_streaming_stream - Query OCI Streaming Service Streams using SQL"
description: "Allows users to query OCI Streaming Service Streams."
---

# Table: oci_streaming_stream - Query OCI Streaming Service Streams using SQL

Oracle Cloud Infrastructure's Streaming Service is a fully managed, scalable, and durable solution for ingesting and consuming high-volume data streams in real time. It allows you to collect, process, and analyze streaming data, such as application logs, telemetry, and other data, in a fault-tolerant way. The service is designed to support streaming use cases, such as analytics, application monitoring, and telemetry, among others.

## Table Usage Guide

The `oci_streaming_stream` table provides insights into streams within Oracle Cloud Infrastructure's Streaming Service. As a data engineer, explore stream-specific details through this table, including partitions, retention periods, and associated metadata. Utilize it to uncover information about streams, such as those with long retention periods, the partition distribution within streams, and the verification of stream properties.

## Examples

### Basic info
Explore which streams in your Oracle Cloud Infrastructure are active and when they were created. This can help you manage and track your resources effectively.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_streaming_stream;
```

### List streams that are not active
Uncover the details of inactive streams within your environment. This can be useful for identifying potential resource waste or areas for optimization.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_streaming_stream
where
  lifecycle_state <> 'ACTIVE';
```

### List streams with retention period greater than 24 hrs
Discover the segments that have a retention period longer than a day. This is useful for identifying and managing streams that require a longer data storage period.

```sql
select
  name,
  id,
  lifecycle_state,
  time_created,
  retention_in_hours
from
  oci_streaming_stream
where
  retention_in_hours > 24;
```