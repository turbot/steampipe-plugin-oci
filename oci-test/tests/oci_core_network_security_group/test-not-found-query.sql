select display_name, id, freeform_tags, lifecycle_state
from oci.oci_core_network_security_group
where id = '{{ output.resource_id.value }}aa';