select db_name, id, lifecycle_state
from oci.oci_database_database
where id = '{{ output.resource_id.value }}';