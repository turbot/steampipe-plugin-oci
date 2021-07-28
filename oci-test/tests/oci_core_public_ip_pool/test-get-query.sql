select display_name, id, lifecycle_state, freeform_tags
from oci.oci_core_public_ip_pool
where id = '{{ output.resource_id.value }}';