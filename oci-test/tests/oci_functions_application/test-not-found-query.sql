select id, display_name, lifecycle_state
from oci.oci_functions_application
where id = '{{ output.resource_id.value }}::dummy';