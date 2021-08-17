select id, name, lifecycle_state
from oci.oci_analytics_instance
where id = 'demo-{{ output.resource_id.value }}';