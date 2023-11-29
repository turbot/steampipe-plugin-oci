---
title: "Steampipe Table: oci_ai_anomaly_detection_project - Query OCI AI Anomaly Detection Projects using SQL"
description: "Allows users to query OCI AI Anomaly Detection Projects."
---

# Table: oci_ai_anomaly_detection_project - Query OCI AI Anomaly Detection Projects using SQL

The OCI AI Anomaly Detection service provides a set of tools within Oracle Cloud Infrastructure that allows you to monitor and respond to anomalies in your data. It uses machine learning and statistical algorithms to detect outliers and unexpected patterns in time-series data. This service is particularly useful in identifying potential issues across your applications and infrastructure.

## Table Usage Guide

The `oci_ai_anomaly_detection_project` table provides insights into anomaly detection projects within Oracle Cloud Infrastructure's AI service. As a data scientist or ML engineer, explore project-specific details through this table, including the project's ID, compartment ID, time created, and lifecycle state. Utilize it to manage and monitor your anomaly detection projects, such as tracking the lifecycle state of projects, understanding the distribution of projects across compartments, and retrieving specific project details.

## Examples

### Basic info
Explore the basic details of anomaly detection projects in your Oracle Cloud Infrastructure. This can help in understanding the state of these projects and their lifecycle.

```sql
select
  id,
  display_name,
  description,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_project;
```

### List the projects created in the last 30 days
Explore which anomaly detection projects have been initiated in the past month. This can help you keep track of recent activities and understand the current focus areas in your organization.

```sql
select
  id,
  display_name,
  description,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_project
where
  time_created >= now() - interval '30' day;
```

### List projects which are not active
Explore which AI anomaly detection projects in your OCI environment are not currently active. This could be useful to identify unused resources or potential areas for cost reduction.

```sql
select
  id,
  display_name,
  description,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_project
where
  lifecycle_state <> 'ACTIVE';
```

### List the models associated to the project
Explore the relationships between various projects and their associated models. This query can be particularly useful in managing and tracking the lifecycle of models within specific projects, providing insights into their creation time and current status.

```sql
select
  p.id as project_id,
  p.display_name as project_name,
  m.id as model_id,
  m.time_created as model_created_time,
  m.display_name as model_name,
  m.lifecycle_state as model_lifecycle_state
from
  oci_ai_anomaly_detection_project as p
  left join oci_ai_anomaly_detection_model as m on p.id = m.project_id;
```

### List the data asset is associated to the project
Explore which data assets are linked to specific projects. This is useful for understanding the distribution and utilization of data assets across different projects.

```sql
select
  p.id as project_id,
  p.display_name as project_name,
  d.id as data_asset_id,
  d.time_created as data_asset_created_time,
  d.display_name as data_asset_name,
  d.lifecycle_state as data_asset_lifecycle_state
from
  oci_ai_anomaly_detection_project as p
  left join oci_ai_anomaly_detection_data_asset as d on p.id = d.project_id;
```