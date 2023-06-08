# Table: oci_ai_anomaly_detection_data_asset

A data asset is an asset related to the anomaly detection project.

## Examples

### Basic info

```sql
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

```sql
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

### List the data assets that are attached to private endpoints

```sql
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

```sql
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

```sql
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