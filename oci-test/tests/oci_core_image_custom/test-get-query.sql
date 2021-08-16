select display_name, id, freeform_tags, lifecycle_state
from oci.oci_core_image_custom
where id = '{{ output.resource_id.value }}';