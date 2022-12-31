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