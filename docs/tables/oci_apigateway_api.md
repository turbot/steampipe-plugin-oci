---
title: "Steampipe Table: oci_apigateway_api - Query OCI API Gateway APIs using SQL"
description: "Allows users to query API Gateway APIs in Oracle Cloud Infrastructure."
---

# Table: oci_apigateway_api - Query OCI API Gateway APIs using SQL

The OCI API Gateway service provides developers with a solution to create, manage, and deploy APIs to expose OCI services or other HTTP backends. It enables the secure and efficient handling of HTTP-based API requests and responses. The service offers features such as request and response transformation, API versioning, and request validation.

## Table Usage Guide

The `oci_apigateway_api` table provides insights into APIs within OCI API Gateway service. As a developer, you can explore API-specific details through this table, including the API's deployment status, lifecycle state, and associated metadata. Utilize it to uncover information about APIs, such as their configurations, the routes they expose, and their overall performance and health.

## Examples

### Basic info
Explore the basic information of your API Gateway in Oracle Cloud Infrastructure to understand its creation time and current lifecycle state. This can be useful in assessing the overall status and management of your APIs.

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_apigateway_api;
```


### List active APIs
Explore which APIs are currently active in your environment. This is beneficial in managing resources and ensuring system efficiency.

```sql
select
  id,
  display_name,
  time_created,
  lifecycle_state as state
from
  oci_apigateway_api
where
  lifecycle_state = 'ACTIVE';
```


### List APIs older than 90 days
Explore which APIs have been created more than 90 days ago. This can be useful for maintaining the health of your system by identifying and potentially updating or removing outdated APIs.

```sql
select
  id,
  lifecycle_state,
  time_created
from
  oci_apigateway_api
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```