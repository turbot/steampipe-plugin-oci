select display_name, id, lifecycle_state, image_type, database_version, patch_set, freeform_tags
from oci.oci_database_software_image
where id = 'demo-test-{{ output.resource_id.value }}';
