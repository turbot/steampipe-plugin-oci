select display_name, id, time_created, lifecycle_state
from oci.oci_mysql_configuration
where id = '{{ output.resource_id.value }}nf';