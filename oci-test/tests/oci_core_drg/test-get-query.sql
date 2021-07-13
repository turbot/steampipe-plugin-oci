select display_name, id, lifecycle_state, freeform_tags
from oci.oci_core_drg
where id = '{{ output.resource_id.value }}';