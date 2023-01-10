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