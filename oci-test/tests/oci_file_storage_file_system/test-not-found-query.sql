select id, display_name, lifecycle_state
from oci.oci_file_storage_file_system
where id = '{{ output.resource_id.value }}::dummy';