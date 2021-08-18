select db_name, id
from oci.oci_database_db_system
where id = '{{ output.resource_id.value }}aa';