# Table: oci_ai_anomaly_detection_model

A model is an asset related to the anomaly detection project.

## Examples

### Basic info

```sql
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

```sql
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

### List models which are not active

```sql
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

```sql
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

