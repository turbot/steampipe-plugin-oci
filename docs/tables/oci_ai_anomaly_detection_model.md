---
title: "Steampipe Table: oci_ai_anomaly_detection_model - Query OCI AI Anomaly Detection Models using SQL"
description: "Allows users to query OCI AI Anomaly Detection Models."
---

# Table: oci_ai_anomaly_detection_model - Query OCI AI Anomaly Detection Models using SQL

OCI AI Anomaly Detection is a service within Oracle Cloud Infrastructure that enables users to detect anomalies in their data by using machine learning. It allows users to build and manage models that can be used for anomaly detection, providing a way to easily identify outliers in complex datasets. This service is particularly useful for applications such as fraud detection, fault detection, and operational optimization.

## Table Usage Guide

The `oci_ai_anomaly_detection_model` table provides insights into AI anomaly detection models within Oracle Cloud Infrastructure. As a data scientist or AI engineer, you can explore model-specific details through this table, including model configurations, training data summaries, and associated metadata. Use it to uncover information about models, such as their status, the type of algorithm used, and the time of the last update.

## Examples

### Basic info
Explore which anomaly detection models are currently active within your project. This allows you to assess their lifecycle states and details, providing insights into the performance and health of each model.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model;
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model;
```

### List the models created in the last 30 days
Uncover the details of recently created AI anomaly detection models within the last month. This is useful for keeping track of new additions and understanding their lifecycle details and states.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model
where
  time_created >= datetime('now', '-30 day');
```

### List models which are not active
Determine the areas in which AI anomaly detection models are not currently active. This can be useful for identifying models that may be outdated or no longer in use, enabling more efficient resource management.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  lifecycle_details,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_model
where
  lifecycle_state <> 'ACTIVE';
```

### List the project details the model is associated to
Discover the segments that link a specific model to its associated projects in the AI anomaly detection system. This allows for a comprehensive overview of project details, aiding in project management and anomaly detection model utilization.

```sql+postgres
select
  m.id as data_asset_id,
  m.display_name as data_asset_name,
  p.id as project_id,
  p.time_created as project_created_time,
  p.display_name as project_name,
  p.lifecycle_state as project_lifecycle_state
from
  oci_ai_anomaly_detection_model as m
  left join oci_ai_anomaly_detection_project as p on m.project_id = p.id;
```

```sql+sqlite
select
  m.id as data_asset_id,
  m.display_name as data_asset_name,
  p.id as project_id,
  p.time_created as project_created_time,
  p.display_name as project_name,
  p.lifecycle_state as project_lifecycle_state
from
  oci_ai_anomaly_detection_model as m
  left join oci_ai_anomaly_detection_project as p on m.project_id = p.id;
```