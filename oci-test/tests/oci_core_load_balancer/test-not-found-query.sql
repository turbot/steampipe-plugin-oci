select display_name, id, lifecycle_state
from oci.oci_core_load_balancer
where id = '{{ output.resource_id.value }}aa';