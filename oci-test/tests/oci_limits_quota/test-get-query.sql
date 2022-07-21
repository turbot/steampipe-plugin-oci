select name, id, lifecycle_state
from oci.oci_limits_quota
where id = '{{ output.resource_id.value }}';