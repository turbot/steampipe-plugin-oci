select display_name, id, lifecycle_state
from oci.oci_core_image
where id = '{{ output.resource_id.value }}';