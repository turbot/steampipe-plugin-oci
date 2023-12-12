---
title: "Steampipe Table: oci_ai_anomaly_detection_data_asset - Query OCI AI Anomaly Detection Data Assets using SQL"
description: "Allows users to query OCI AI Anomaly Detection Data Assets."
---

# Table: oci_ai_anomaly_detection_data_asset - Query OCI AI Anomaly Detection Data Assets using SQL

Oracle Cloud Infrastructure's AI Anomaly Detection service enables you to build and deploy machine learning models to detect anomalies in your time series data. The AI Anomaly Detection Data Asset is a resource in this service that represents a data source from which the service can ingest data for training and inference. It supports data sources such as Object Storage and Data Flow.

## Table Usage Guide

The `oci_ai_anomaly_detection_data_asset` table provides insights into the Data Assets within OCI AI Anomaly Detection service. As a Data Scientist or AI Engineer, you can explore data asset-specific details through this table, including data source type, data source details, and associated metadata. Utilize it to uncover information about your data assets, such as the data source they are linked to and the type of data source, assisting you in managing and organizing your AI Anomaly Detection resources.

## Examples

### Basic info
Explore the basic details of your AI anomaly detection data assets in Oracle Cloud Infrastructure to understand their states and associated projects. This can be useful for auditing and managing your AI resources effectively.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset;
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset;
```

### List the data assets created in the last 30 days
Discover the recently created data assets in the last 30 days to keep track of your project's growth and evolution. This allows for timely assessment and review of newly added resources in your project.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  time_created >= datetime('now','-30 day');
```

### List the data assets that are attached to private endpoints
Determine the areas in which data assets are linked to private endpoints. This is useful for understanding the security and accessibility of your data assets within your project.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  private_endpoint_id is not null;
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  private_endpoint_id is not null;
```

### List data assets which are not active
Explore which data assets in your AI anomaly detection project are not currently active. This can help manage resources and identify any assets that may need to be reactivated or removed.

```sql+postgres
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  lifecycle_state <> 'ACTIVE';
```

```sql+sqlite
select
  id,
  display_name,
  project_id,
  description,
  private_endpoint_id,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_data_asset
where
  lifecycle_state <> 'ACTIVE';
```

### List the project details the data asset is associated to
Discover the connections between data assets and their associated projects. This query can be used to gain insights into the lifecycle of different projects and their corresponding data assets, providing a clear overview of project timelines and their related data assets.

```sql+postgres
select
  d.id as data_asset_id,
  d.display_name as data_asset_name,
  p.id as project_id,
  p.time_created as project_created_time,
  p.display_name as project_name,
  p.lifecycle_state as project_lifecycle_state
from
  oci_ai_anomaly_detection_data_asset as d
  left join oci_ai_anomaly_detection_project as p on d.project_id = p.id;
```

```sql+sqlite
select
  d.id as data_asset_id,
  d.display_name as data_asset_name,
  p.id as project_id,
  p.time_created as project_created_time,
  p.display_name as project_name,
  p.lifecycle_state as project_lifecycle_state
from
  oci_ai_anomaly_detection_data_asset as d
  left join oci_ai_anomaly_detection_project as p on d.project_id = p.id;
```