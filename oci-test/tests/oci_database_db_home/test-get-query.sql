select display_name, id, lifecycle_state
from oci.oci_database_db_home
where id = '{{ output.resource_id.value }}';