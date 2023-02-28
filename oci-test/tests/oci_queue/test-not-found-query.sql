select id, display_name, lifecycle_state
from oci.oci_queue
where id = '{{ output.resource_id.value }}::dummy';