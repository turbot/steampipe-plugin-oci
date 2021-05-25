select id, display_name, lifecycle_state, availability_domain, is_clone_parent, is_hydrated, freeform_tags
from oci.oci_file_storage_file_system
where id = '{{ output.resource_id.value }}';