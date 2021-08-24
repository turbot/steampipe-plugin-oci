select id, display_name, availability_domain, lifecycle_state, subnet_id, freeform_tags
from oci.oci_file_storage_mount_target
where display_name = '{{ output.resource_name.value }}';