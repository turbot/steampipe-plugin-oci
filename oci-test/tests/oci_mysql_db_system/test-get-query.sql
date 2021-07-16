select display_name, id, lifecycle_state
from oci.oci_mysql_db_system
where id = '{{ output.resource_id.value }}';