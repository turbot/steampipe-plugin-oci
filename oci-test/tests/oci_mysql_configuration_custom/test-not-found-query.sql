select display_name, id, time_created, lifecycle_state
from oci.oci_mysql_configuration_custom
where id = 'demo-{{ output.resource_id.value }}';
