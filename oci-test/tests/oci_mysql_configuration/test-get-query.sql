select display_name, id, lifecycle_state, shape_name, type
from oci.oci_mysql_configuration
where id = '{{ output.resource_id.value }}';