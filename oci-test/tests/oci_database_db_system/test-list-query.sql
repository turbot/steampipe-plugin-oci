select db_name, id, lifecycle_state
from oci.oci_database_db_system
where db_name = '{{ output.db_name.value }}';