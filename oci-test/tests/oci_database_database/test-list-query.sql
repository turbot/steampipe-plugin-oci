select db_name, id, lifecycle_state
from oci.oci_database_database
where db_name = '{{ output.db_name.value }}';