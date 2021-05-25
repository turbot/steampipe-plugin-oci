select id, name, lifecycle_state
from oci.oci_file_storage_snapshot
where id = '{{ output.resource_id.value }}::dummy';