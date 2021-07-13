select name, id, lifecycle_state, description, inactive_status
from oci.oci_identity_network_source
where id = '{{ output.resource_id.value }}';