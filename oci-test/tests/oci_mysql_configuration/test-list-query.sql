select display_name, id, lifecycle_state, shape_name, type
from oci.oci_mysql_configuration
where display_name = '{{ resourceName }}';