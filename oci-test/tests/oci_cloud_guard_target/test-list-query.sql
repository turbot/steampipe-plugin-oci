select id, name, lifecycle_state
from oci.oci_cloud_guard_target
where name = '{{ output.resource_name.value }}';