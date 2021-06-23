select id, name, lifecycle_state, is_clone_source, freeform_tags
from oci.oci_file_storage_snapshot
where id = '{{ output.resource_id.value }}';