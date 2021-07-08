select db_name, id, lifecycle_state
from oci.oci_database_autonomous_database
where id = '{{ output.resource_id.value }}';