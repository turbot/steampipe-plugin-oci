select id, name, lifecycle_state
from oci.oci_cloud_guard_target
where id = '{{ output.resource_id.value }}::dummy';