select display_name, id, lifecycle_state, is_private, shape_name
from oci.oci_core_load_balancer
where id = '{{ output.resource_id.value }}';