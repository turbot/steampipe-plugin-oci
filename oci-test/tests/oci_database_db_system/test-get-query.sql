select db_name, id, lifecycle_state
from oci.oci_database_db_system
where id = '{{ output.resource_id.value }}';