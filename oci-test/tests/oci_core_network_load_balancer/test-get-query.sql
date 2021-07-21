select display_name, id, subnet_id, lifecycle_state, freeform_tags, is_private
from oci.oci_core_network_load_balancer
where id = '{{ output.resource_id.value }}';