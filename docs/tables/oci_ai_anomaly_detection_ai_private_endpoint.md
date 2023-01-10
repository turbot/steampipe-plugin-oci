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