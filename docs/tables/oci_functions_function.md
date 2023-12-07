---
title: "Steampipe Table: oci_functions_function - Query OCI Functions using SQL"
description: "Allows users to query OCI Functions."
---

# Table: oci_functions_function - Query OCI Functions using SQL

Oracle Cloud Infrastructure (OCI) Functions is a fully managed, multi-tenant, highly scalable, on-demand, Functions-as-a-Service (FaaS) platform. It is powered by the open source Fn Project and allows developers to write and deploy code without provisioning and managing servers. OCI Functions can be written in Java, Python, Node, Go, and Ruby, and they can be triggered by HTTP requests or OCI Events Service.

## Table Usage Guide

The `oci_functions_function` table provides insights into functions within Oracle Cloud Infrastructure (OCI) Functions. As a software developer or DevOps engineer, you can explore function-specific details through this table, including the runtime, memory limits, and associated metadata. Use it to uncover information about functions, such as their execution environment, the maximum amount of time that the function is allowed to run, and the state of the function.

## Examples

### Basic info
Explore the basic information of your Oracle Cloud Infrastructure functions to gain insights into details such as their lifecycle state and associated application ID. This can help in understanding the current state and configuration of your serverless applications.

```sql+postgres
select
  display_name,
  id,
  application_id,
  lifecycle_state,
  image,
  image_digest
from
  oci_functions_function;
```

```sql+sqlite
select
  display_name,
  id,
  application_id,
  lifecycle_state,
  image,
  image_digest
from
  oci_functions_function;
```


### List functions where trace configuration is disabled
Identify functions where trace configuration is not enabled to understand potential areas of risk or lack of visibility in your application performance monitoring.

```sql+postgres
select
  display_name,
  id,
  application_id,
  trace_config -> 'isEnabled' as trace_config_is_enabled
from
  oci_functions_function
where
  not (trace_config -> 'isEnabled') :: bool;
```

```sql+sqlite
select
  display_name,
  id,
  application_id,
  json_extract(trace_config, '$.isEnabled') as trace_config_is_enabled
from
  oci_functions_function
where
  not json_extract(trace_config, '$.isEnabled');
```


### List functions where memory is greater than 100 MB
Discover the functions that are using more than 100 MB of memory. This can be useful to identify high-memory usage functions that may need optimization or resource allocation adjustments.

```sql+postgres
select
  display_name,
  id,
  application_id,
  memory_in_mbs
from
  oci_functions_function
where
  memory_in_mbs > 100;
```

```sql+sqlite
select
  display_name,
  id,
  application_id,
  memory_in_mbs
from
  oci_functions_function
where
  memory_in_mbs > 100;
```