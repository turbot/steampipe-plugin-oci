select id, display_name, lifecycle_state, freeform_tags
from oci.oci_queue
where id = '{{ output.resource_id.value }}';