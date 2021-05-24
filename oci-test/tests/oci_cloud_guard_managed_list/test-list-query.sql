select id, name, lifecycle_state
from oci.oci_cloud_guard_managed_list
where name = '{{ output.resource_name.value }}';