select id, name, lifecycle_state, file_system_id
from oci.oci_file_storage_snapshot
where name = '{{ output.resource_name.value }}';