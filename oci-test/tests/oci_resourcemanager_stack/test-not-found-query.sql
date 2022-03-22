select display_name, id, time_created, lifecycle_state
from oci.oci_resourcemanager_stack
where id = '{{ output.resource_id.value }}nf';