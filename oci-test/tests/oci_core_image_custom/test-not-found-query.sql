select display_name, id, freeform_tags, size_in_mbs, lifecycle_state
from oci.oci_core_image_custom
where id = '{{ output.resource_id.value }}nf';