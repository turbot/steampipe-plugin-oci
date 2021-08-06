select id, display_name, availability_domain, lifecycle_state, subnet_id, freeform_tags
from oci.oci_file_storage_mount_target
where id = '{{ output.resource_id.value }}';