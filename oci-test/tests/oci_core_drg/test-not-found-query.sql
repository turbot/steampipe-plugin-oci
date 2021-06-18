select display_name, id, time_created, lifecycle_state
from oci.oci_core_drg
where id = '{{ output.resource_id.value }}dummy';