# Table: oci_ai_anomaly_detection_project

A project is a collection of resources for the anomaly detection service.

## Examples

### Basic info

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
