select display_name, id, lifecycle_state
from oci.oci_mysql_db_system
where display_name = '{{ resourceName }}';