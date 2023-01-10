# Table: oci_big_data_service_instance

Oracle Big Data Service is a fully managed, automated cloud service that provides enterprises with a cost-effective Hadoop environment. Customers easily create secure and scalable Hadoop-based data lakes that can quickly process large amounts of data.

## Examples

### Basic info

```sql
select
    id,
    display_name,
    is_high_availability,
    is_secure,
    is_cloud_sql_configured,
    nodes,
    number_of_nodes,
    cluster_version,
    network_config,
    cluster_details,
    cloud_sql_details,
    created_by,
    bootstrap_script_url,
    kms_key_id,
    cluster_profile,
    lifecycle_state as state
from
    oci_big_data_service_instance;
```