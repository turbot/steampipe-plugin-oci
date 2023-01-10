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