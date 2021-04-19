select name, id, lifecycle_state, title
from oci.oci_logging_log
where id  = '{{ output.resource_id.value }}' and log_group_id = '{{ output.log_group_id.value }}';