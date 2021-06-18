select id, display_name, lifecycle_state, availability_domain
from oci.oci_file_storage_file_system
where display_name = '{{ output.resource_name.value }}';