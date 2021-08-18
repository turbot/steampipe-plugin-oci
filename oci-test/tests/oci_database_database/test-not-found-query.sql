select db_name, id
from oci.oci_database_database
where id = '{{ output.resource_id.value }}aa';