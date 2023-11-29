---
title: "Steampipe Table: oci_queue_queue - Query OCI Queue Service Queues using SQL"
description: "Allows users to query OCI Queue Service Queues."
---

# Table: oci_queue_queue - Query OCI Queue Service Queues using SQL

The Oracle Cloud Infrastructure (OCI) Queue Service is a fully managed, scalable, and durable messaging service for asynchronous communication among microservices, distributed systems, and serverless applications. It allows you to send, store, and receive messages between software components at any volume, without losing messages or requiring other services to be always available. OCI Queue Service is built on OCI Streaming and provides a simple, RESTful API for messaging.

## Table Usage Guide

The `oci_queue_queue` table provides insights into the queues within the OCI Queue Service. As a systems engineer, you can explore queue-specific details through this table, including queue type, retention period, and associated metadata. This table is particularly useful for managing and monitoring the performance and health of your queues, ensuring efficient communication between your microservices and applications.

## Examples

### Basic info
Explore which queues are active and their associated details to better manage and monitor your Oracle Cloud Infrastructure (OCI) resources. This can help you maintain an efficient and organized OCI environment.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  lifecycle_state
from
  oci_queue_queue;
```

### List queues not in the active state
Determine the areas in which queues are not actively functioning. This can be useful for identifying potential bottlenecks or interruptions in data flow.

```sql
select
  display_name,
  id,
  lifecycle_state as state
from
  oci_queue_queue
where
  lifecycle_state <> 'ACTIVE';
```

### List queues that are not encrypted
Discover the segments that are not encrypted in your system. This query is useful for identifying potential security risks and ensuring that all queues in your system are properly protected.

```sql
select
  display_name,
  id,
  lifecycle_state as state,
  compartment_id,
  region
from
  oci_queue_queue
where 
  custom_encryption_key_id is null;
```  

### Get details of each queue
Explore the characteristics of each queue, such as retention and visibility duration, to manage and optimize your queue configurations effectively. This allows you to understand queue behaviors, identify potential bottlenecks, and implement necessary adjustments for improved performance.

```sql
select
  display_name,
  retention_in_seconds,
  visibility_in_seconds,
  timeout_in_seconds,
  messages_endpoint,
  dead_letter_queue_delivery_count,
  defined_tags,
  id
from
  oci_queue_queue;
```