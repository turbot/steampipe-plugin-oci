select display_name, id, vcn_id, lifecycle_state, freeform_tags
from oci.oci_core_network_security_group
where id = '{{ output.resource_id.value }}';