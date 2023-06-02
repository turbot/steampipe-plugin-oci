# Table: oci_ai_anomaly_detection_ai_private_endpoint

A private network reverse connection creates a connection from service to customer subnet over a private network.

## Examples

### Basic info

```sql
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint;
```

### List the AI private endpoints created in the last 30 days

```sql
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  time_created >= now() - interval '30' day;
```

### List the AI private endpoints that have attached data assets

```sql
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  attached_data_assets is not null;
```

### List AI private endpoints which are not active

```sql
select
  id,
  subnet_id,
  dns_zones,
  display_name,
  lifecycle_details,
  attached_data_assets,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  lifecycle_state <> 'ACTIVE';
```

### List DNS zones associated to the AI private endpoints

```sql
select
  id,
  subnet_id,
  jsonb_pretty(dns_zones) as dns_zones,
  display_name,
  lifecycle_state as state
from
  oci_ai_anomaly_detection_ai_private_endpoint
where
  dns_zones is not null;
```
