select id, display_name, lifecycle_state
from oci.oci_queue
where display_name = '{{ output.display_name.value }}';