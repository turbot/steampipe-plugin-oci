select id, display_name, lifecycle_state
from oci.oci_functions_application
where display_name = '{{ output.display_name.value }}';