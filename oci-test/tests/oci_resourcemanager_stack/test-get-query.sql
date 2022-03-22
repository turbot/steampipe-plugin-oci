select display_name, id, lifecycle_state, time_created
from oci.oci_resourcemanager_stack
where id = '{{ output.resource_id.value }}';