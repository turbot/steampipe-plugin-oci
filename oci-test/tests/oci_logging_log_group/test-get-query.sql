select display_name, id, lifecycle_state, description, title
from oci.oci_logging_log_group
where id  = '{{ output.resource_id.value }}';