select display_name, id, freeform_tags, lifecycle_state, dns_label
from oci.oci_core_subnet
where id = '{{ output.resource_id.value }}';